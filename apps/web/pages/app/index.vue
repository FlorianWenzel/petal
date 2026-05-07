<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import { useJournal } from '~/composables/useJournal'
import {
  averageColor,
  colorFor,
  daysInMonth,
  type Mood,
} from '~/lib/journal'

definePageMeta({
  middleware: 'auth',
  layout: false,
})
useSeoMeta({ title: 'my flower — petal', robots: 'noindex' })

const { user, logout, renameUser } = useAuth()

const editingName = ref(false)
const nameDraft = ref('')
const nameError = ref<string | null>(null)
const nameSaving = ref(false)

function startRename() {
  if (!user.value) return
  nameDraft.value = user.value.username
  nameError.value = null
  editingName.value = true
}

function cancelRename() {
  editingName.value = false
  nameError.value = null
}

async function commitRename() {
  if (!user.value) return
  const trimmed = nameDraft.value.trim().toLowerCase()
  if (trimmed === user.value.username) {
    editingName.value = false
    return
  }
  nameSaving.value = true
  nameError.value = null
  try {
    await renameUser(trimmed)
    editingName.value = false
  } catch (e: any) {
    nameError.value = e?.data?.message ?? 'could not rename'
  } finally {
    nameSaving.value = false
  }
}

const today = new Date()
const todayY = today.getFullYear()
const todayM = today.getMonth()
const todayD = today.getDate()

const viewYear = ref(todayY)
const viewMonth = ref(todayM)

const { entries, load, save } = useJournal()

const dialogOpen = ref(false)
const dialogDay = ref(1)

const monthDays = computed(() => daysInMonth(viewYear.value, viewMonth.value))

const monthKey = computed(
  () => `${viewYear.value}-${String(viewMonth.value + 1).padStart(2, '0')}`,
)

const entriesByDay = computed<Record<number, string>>(() => {
  const out: Record<number, string> = {}
  for (const [k, v] of Object.entries(entries.value)) {
    out[Number(k)] = colorFor(v.mood)
  }
  return out
})

const todayDay = computed(() =>
  viewYear.value === todayY && viewMonth.value === todayM ? todayD : null,
)

const isCurrentMonth = computed(
  () => viewYear.value === todayY && viewMonth.value === todayM,
)

const monthLabel = computed(() =>
  new Date(viewYear.value, viewMonth.value, 1).toLocaleString(undefined, {
    month: 'long',
    year: 'numeric',
  }),
)

const filledCount = computed(() => Object.keys(entries.value).length)

const centerColor = computed(() => averageColor(entries.value))

const initialMood = computed<Mood | null>(
  () => entries.value[String(dialogDay.value)]?.mood ?? null,
)

const initialText = computed(
  () => entries.value[String(dialogDay.value)]?.text ?? '',
)

function prev() {
  if (viewMonth.value === 0) {
    viewMonth.value = 11
    viewYear.value -= 1
  } else {
    viewMonth.value -= 1
  }
}

function next() {
  if (isCurrentMonth.value) return
  if (viewMonth.value === 11) {
    viewMonth.value = 0
    viewYear.value += 1
  } else {
    viewMonth.value += 1
  }
}

function onPetalClick(day: number) {
  if (viewYear.value === todayY && viewMonth.value === todayM && day > todayD) return
  dialogDay.value = day
  dialogOpen.value = true
}

async function onSave(payload: { mood: Mood; text: string }) {
  await save(viewYear.value, viewMonth.value, dialogDay.value, payload.mood, payload.text)
  dialogOpen.value = false
}

watch([viewYear, viewMonth], () => load(viewYear.value, viewMonth.value), { immediate: true })
</script>

<template>
  <div class="flex flex-col h-screen w-full">
    <div class="shrink-0 px-4 pt-3">
      <div class="max-w-3xl mx-auto flex items-center justify-between gap-4 text-sm">
        <div class="min-w-0 flex-1">
          <template v-if="editingName">
            <form class="flex items-center gap-2" @submit.prevent="commitRename">
              <input
                v-model="nameDraft"
                autofocus
                maxlength="24"
                class="px-2 py-1 rounded-md bg-white/80 border border-stone-300 text-stone-800 font-mono w-44 focus:outline-none focus:ring-2 focus:ring-stone-300"
                :disabled="nameSaving"
                @keydown.esc="cancelRename"
              />
              <button
                type="submit"
                :disabled="nameSaving"
                class="text-xs px-2 py-1 rounded-md bg-stone-700 text-white hover:bg-stone-800 disabled:opacity-50"
              >
                save
              </button>
              <button
                type="button"
                class="text-xs px-2 py-1 rounded-md text-stone-500 hover:text-stone-700"
                @click="cancelRename"
              >
                cancel
              </button>
            </form>
            <div v-if="nameError" class="text-xs text-red-600 mt-1">{{ nameError }}</div>
          </template>
          <template v-else>
            <button
              class="group inline-flex items-center gap-2 text-stone-600 hover:text-stone-900 transition"
              :title="'rename'"
              @click="startRename"
            >
              <span class="font-mono">{{ user?.username ?? '…' }}</span>
              <span class="text-xs opacity-0 group-hover:opacity-100 transition">rename</span>
            </button>
          </template>
        </div>
        <button
          class="shrink-0 text-xs px-3 py-1.5 rounded-md bg-white/70 border border-white text-stone-600 hover:text-stone-900 hover:bg-white transition"
          @click="logout"
        >
          log out
        </button>
      </div>
    </div>
    <div class="shrink-0 px-4 pt-4">
      <div class="max-w-3xl mx-auto flex items-center justify-center gap-4">
        <button
          class="w-10 h-10 rounded-full bg-white/70 border border-white text-stone-700 hover:bg-white transition flex items-center justify-center text-xl"
          aria-label="previous month"
          @click="prev"
        >
          ‹
        </button>
        <div class="text-center min-w-[12rem]">
          <div class="text-xl font-semibold text-stone-800">{{ monthLabel }}</div>
          <div class="text-xs text-stone-500 font-mono">
            {{ filledCount }} / {{ monthDays }} days
          </div>
        </div>
        <button
          :disabled="isCurrentMonth"
          class="w-10 h-10 rounded-full bg-white/70 border border-white text-stone-700 hover:bg-white transition flex items-center justify-center text-xl disabled:opacity-30 disabled:cursor-not-allowed"
          aria-label="next month"
          @click="next"
        >
          ›
        </button>
      </div>
    </div>

    <div class="flex-1 relative">
      <ClientOnly>
        <Flower3D
          :days-in-month="monthDays"
          :entries-by-day="entriesByDay"
          :today-day="todayDay"
          :month-key="monthKey"
          :center-color="centerColor"
          @petal-click="onPetalClick"
        />
      </ClientOnly>
    </div>

    <div class="shrink-0 px-4 py-3 text-center text-sm text-stone-500">
      tap a petal to log how the day felt
    </div>

    <ClientOnly>
      <EntryDialog
        :open="dialogOpen"
        :year="viewYear"
        :month="viewMonth"
        :day="dialogDay"
        :initial-mood="initialMood"
        :initial-text="initialText"
        @save="onSave"
        @close="dialogOpen = false"
      />
    </ClientOnly>
  </div>
</template>

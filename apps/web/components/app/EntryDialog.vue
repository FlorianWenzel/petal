<script setup lang="ts">
import { ref, watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import { BubbleMenu } from '@tiptap/vue-3/menus'
import StarterKit from '@tiptap/starter-kit'
import { MOODS, type Mood } from '~/lib/journal'

const props = defineProps<{
  open: boolean
  year: number
  month: number
  day: number
  initialMood: Mood | null
  initialText: string
}>()

const emit = defineEmits<{
  save: [payload: { mood: Mood; text: string }]
  close: []
}>()

const selected = ref<Mood | null>(props.initialMood)

const editor = useEditor({
  extensions: [
    StarterKit.configure({
      heading: { levels: [2, 3] },
    }),
  ],
  content: props.initialText,
  editorProps: {
    attributes: {
      class:
        'prose-petal min-h-[8rem] max-h-[14rem] overflow-y-auto px-3 py-3 rounded-xl bg-stone-50 border border-stone-200 focus:outline-none focus:border-stone-400 transition',
    },
  },
})

watch(
  () => [props.open, props.initialMood, props.initialText] as const,
  ([open, mood, text]) => {
    if (open) {
      selected.value = mood
      if (editor.value && editor.value.getHTML() !== text) {
        editor.value.commands.setContent(text || '')
      }
    }
  },
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})

const monthName = (y: number, m: number) =>
  new Date(y, m, 1).toLocaleString(undefined, { month: 'long', year: 'numeric' })

function save() {
  if (!selected.value) return
  emit('save', { mood: selected.value, text: editor.value?.getHTML() ?? '' })
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-stone-900/30 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-w-lg rounded-2xl bg-white shadow-2xl border border-white/80 p-6 flex flex-col gap-4"
      >
        <div class="flex items-baseline justify-between">
          <h2 class="text-xl font-semibold text-stone-800">
            {{ day }} {{ monthName(year, month) }}
          </h2>
          <button
            @click="emit('close')"
            class="text-stone-400 hover:text-stone-700 text-xl leading-none"
            aria-label="close"
          >
            ×
          </button>
        </div>

        <div class="flex flex-col gap-2">
          <p class="text-sm text-stone-500">how did this day feel?</p>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="m in MOODS"
              :key="m.id"
              @click="selected = m.id"
              :class="[
                'px-4 py-2 rounded-full text-sm font-medium transition-all border-2',
                selected === m.id
                  ? 'border-stone-700 scale-105 shadow'
                  : 'border-transparent hover:scale-105',
              ]"
              :style="{ backgroundColor: m.color, color: '#3a2f3f' }"
            >
              {{ m.label }}
            </button>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <p class="text-sm text-stone-500">a few words?</p>
          <BubbleMenu
            v-if="editor"
            :editor="editor"
            class="flex gap-1 px-1.5 py-1 rounded-full bg-rose-brand text-rose-brand-ink shadow-lg text-xs"
          >
            <button
              @click="editor.chain().focus().toggleBold().run()"
              :class="[
                'px-2 py-1 rounded-full transition',
                editor.isActive('bold') ? 'bg-white/20' : 'hover:bg-white/10',
              ]"
            >
              <span class="font-bold">B</span>
            </button>
            <button
              @click="editor.chain().focus().toggleItalic().run()"
              :class="[
                'px-2 py-1 rounded-full transition',
                editor.isActive('italic') ? 'bg-white/20' : 'hover:bg-white/10',
              ]"
            >
              <span class="italic">I</span>
            </button>
            <button
              @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
              :class="[
                'px-2 py-1 rounded-full transition',
                editor.isActive('heading', { level: 2 }) ? 'bg-white/20' : 'hover:bg-white/10',
              ]"
            >
              H
            </button>
            <button
              @click="editor.chain().focus().toggleBulletList().run()"
              :class="[
                'px-2 py-1 rounded-full transition',
                editor.isActive('bulletList') ? 'bg-white/20' : 'hover:bg-white/10',
              ]"
            >
              •
            </button>
            <button
              @click="editor.chain().focus().toggleBlockquote().run()"
              :class="[
                'px-2 py-1 rounded-full transition',
                editor.isActive('blockquote') ? 'bg-white/20' : 'hover:bg-white/10',
              ]"
            >
              ❝
            </button>
          </BubbleMenu>
          <EditorContent :editor="editor" />
        </div>

        <div class="flex justify-end gap-2 pt-1">
          <button
            @click="emit('close')"
            class="px-4 py-2 rounded-full text-sm text-stone-600 hover:bg-stone-100 transition"
          >
            cancel
          </button>
          <button
            @click="save"
            :disabled="!selected"
            class="px-5 py-2 rounded-full font-semibold shadow-md hover:shadow-lg active:scale-95 transition-all disabled:opacity-40 disabled:cursor-not-allowed disabled:active:scale-100"
            :style="{
              backgroundColor: selected ? MOODS.find((m) => m.id === selected)?.color : '#e7e5e4',
              color: '#3a2f3f',
            }"
          >
            save
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style>
.prose-petal {
  font-family: ui-serif, Georgia, 'Times New Roman', serif;
  color: #44403c;
  font-size: 0.95rem;
  line-height: 1.55;
}
.prose-petal p {
  margin: 0 0 0.5em 0;
}
.prose-petal h2 {
  font-size: 1.15rem;
  font-weight: 600;
  margin: 0.5em 0 0.3em 0;
  color: #292524;
}
.prose-petal h3 {
  font-size: 1.05rem;
  font-weight: 600;
  margin: 0.4em 0 0.2em 0;
}
.prose-petal ul {
  list-style: disc;
  padding-left: 1.25rem;
  margin: 0.25em 0;
}
.prose-petal ol {
  list-style: decimal;
  padding-left: 1.25rem;
  margin: 0.25em 0;
}
.prose-petal blockquote {
  border-left: 3px solid #d6d3d1;
  padding-left: 0.75rem;
  color: #57534e;
  font-style: italic;
  margin: 0.5em 0;
}
.prose-petal code {
  background: #f5f5f4;
  padding: 0.1em 0.35em;
  border-radius: 0.25rem;
  font-size: 0.9em;
}
.prose-petal:focus,
.prose-petal *:focus {
  outline: none;
}
.prose-petal p.is-editor-empty:first-child::before {
  content: 'write a thought…';
  color: #a8a29e;
  pointer-events: none;
  height: 0;
  float: left;
}
</style>

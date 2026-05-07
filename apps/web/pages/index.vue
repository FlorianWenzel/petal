<script setup lang="ts">
import { averageColor, colorFor, MOODS, type Mood, type MonthEntries } from '~/lib/journal'

useSeoMeta({
  title: 'petal — a quiet mood journal',
  description:
    'a calm, beautiful mood journal. one petal a day. by the end of the month, your feelings have become a flower.',
  ogTitle: 'petal — a quiet mood journal',
  ogDescription:
    'one petal a day. by the end of the month, your feelings have become a flower.',
  ogType: 'website',
  twitterCard: 'summary_large_image',
})

const demoMonth: MonthEntries = (() => {
  const moods: Mood[] = ['bright', 'tender', 'calm', 'heavy', 'stormy']
  const pattern: Mood[] = [
    'tender', 'calm', 'bright', 'tender', 'bright',
    'calm', 'heavy', 'tender', 'bright', 'calm',
    'stormy', 'tender', 'bright', 'calm', 'tender',
    'heavy', 'calm', 'bright', 'bright', 'tender',
    'bright', 'calm', 'tender', 'heavy', 'bright',
    'calm', 'tender', 'bright', 'bright', 'tender',
  ]
  const out: MonthEntries = {}
  for (let i = 0; i < 30; i++) {
    out[String(i + 1)] = { mood: pattern[i] ?? moods[i % moods.length] }
  }
  return out
})()

const demoEntriesByDay = computed<Record<number, string>>(() => {
  const out: Record<number, string> = {}
  for (const [k, v] of Object.entries(demoMonth)) {
    out[Number(k)] = colorFor(v.mood)
  }
  return out
})

const demoCenterColor = computed(() => averageColor(demoMonth))

const steps = [
  {
    color: '#FFD93D',
    title: 'pick a mood',
    body: 'five soft choices. bright, tender, calm, heavy, stormy. no scales, no scoring.',
  },
  {
    color: '#FF9EBB',
    title: 'write a sentence',
    body: 'one line about the day. or none — the petal alone is already a kind of writing.',
  },
  {
    color: '#B5C7F7',
    title: 'watch it bloom',
    body: 'the flower fills out as the month goes on. the center holds the average of how it all felt.',
  },
]

const quotes = [
  {
    text: 'i don\'t trust apps that try to make me happier. this one barely tries. perfect.',
    name: 'mikael',
    role: 'deeply suspicious of joy',
  },
  {
    text: 'i was going to make a habit tracker. then i drew a flower instead. honestly i don\'t know why you\'d use this either.',
    name: 'flo',
    role: 'made this, allegedly',
  },
  {
    text: 'I am just disappointed this app doesnt feature coconuts at all',
    name: 'nic',
    role: 'journaling expert',
  },
]

function buildMonth(pattern: Mood[], days: number): MonthEntries {
  const out: MonthEntries = {}
  for (let i = 0; i < days; i++) {
    out[String(i + 1)] = { mood: pattern[i % pattern.length] }
  }
  return out
}

function entriesByDayFor(m: MonthEntries): Record<number, string> {
  const out: Record<number, string> = {}
  for (const [k, v] of Object.entries(m)) out[Number(k)] = colorFor(v.mood)
  return out
}

const gallery = (() => {
  const tender = buildMonth(
    [
      'tender', 'tender', 'bright', 'tender', 'calm',
      'tender', 'bright', 'tender', 'tender', 'calm',
      'bright', 'tender', 'tender', 'calm', 'tender',
      'bright', 'tender', 'tender', 'bright', 'calm',
      'tender', 'bright', 'tender', 'calm', 'tender',
      'tender', 'bright', 'tender', 'calm', 'tender', 'bright',
    ],
    31,
  )
  const move = buildMonth(
    [
      'heavy', 'heavy', 'stormy', 'calm', 'heavy',
      'stormy', 'heavy', 'calm', 'heavy', 'stormy',
      'heavy', 'heavy', 'calm', 'stormy', 'heavy',
      'heavy', 'stormy', 'calm', 'heavy', 'heavy',
      'stormy', 'heavy', 'calm', 'stormy', 'heavy',
      'heavy', 'stormy', 'calm',
    ],
    28,
  )
  const finals = buildMonth(
    [
      'stormy', 'bright', 'stormy', 'heavy', 'stormy',
      'bright', 'stormy', 'stormy', 'heavy', 'bright',
      'stormy', 'stormy', 'heavy', 'stormy', 'bright',
      'stormy', 'heavy', 'stormy', 'bright', 'stormy',
      'heavy', 'stormy', 'stormy', 'bright', 'heavy',
      'stormy', 'bright', 'stormy', 'stormy', 'bright',
    ],
    30,
  )
  return [
    { title: 'a tender march', sub: '31 entries · soft pinks', month: tender },
    { title: 'the move', sub: '28 entries · long week', month: move },
    { title: 'studying for finals', sub: '30 entries · burnt edges', month: finals },
  ].map((g) => ({
    ...g,
    entries: entriesByDayFor(g.month),
    center: averageColor(g.month),
    days: Object.keys(g.month).length,
  }))
})()
</script>

<template>
  <main>
    <!-- hero -->
    <section class="px-6 pt-12 pb-24 md:pt-20 md:pb-32">
      <div class="max-w-6xl mx-auto grid md:grid-cols-2 gap-12 md:gap-16 items-center">
        <div class="text-center md:text-left">
          <span
            class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-white/60 border border-white text-xs text-stone-600 font-mono tracking-wide"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-pink-300" />
            a small daily ritual
          </span>
          <h1
            class="mt-6 text-5xl md:text-7xl font-semibold text-stone-800 tracking-tight leading-[1.05]"
          >
            a quiet place<br />for how you feel.
          </h1>
          <p class="mt-6 text-lg md:text-xl text-stone-600 max-w-md mx-auto md:mx-0 leading-relaxed">
            one petal a day. by the end of the month,
            your feelings have become a flower.
          </p>
          <div class="mt-10">
            <NuxtLink
              to="/signup"
              data-petals class="inline-flex items-center gap-2 px-7 py-3.5 rounded-full bg-rose-brand text-rose-brand-ink text-base hover:bg-rose-brand-hover transition shadow-sm"
            >
              start journaling
              <span aria-hidden="true">→</span>
            </NuxtLink>
            <p class="mt-4 text-sm text-stone-500">
              free for the first month. no card needed.
            </p>
          </div>
        </div>

        <div class="flex justify-center md:justify-end">
          <div class="relative w-[480px] max-w-full aspect-square">
            <div
              class="absolute bottom-2 left-1/2 -translate-x-1/2 w-3/4 h-6 rounded-full bg-stone-900/10 blur-2xl"
            />
            <ClientOnly>
              <Flower3D
                preview
                :days-in-month="30"
                :entries-by-day="demoEntriesByDay"
                :today-day="null"
                month-key="hero-demo"
                :center-color="demoCenterColor"
                class="w-full h-full"
              />
              <template #fallback>
                <FlowerPlaceholder :size="480" label="loading flower…" />
              </template>
            </ClientOnly>
          </div>
        </div>
      </div>
    </section>

    <!-- how it works -->
    <section class="px-6 py-20">
      <div class="max-w-5xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-semibold text-stone-800 text-center">
          three small things, once a day.
        </h2>
        <p class="mt-3 text-stone-600 text-center max-w-xl mx-auto">
          it takes about ten seconds. that is the point.
        </p>
        <div class="mt-14 grid md:grid-cols-3 gap-6">
          <div
            v-for="(s, i) in steps"
            :key="s.title"
            class="bg-white/60 rounded-3xl p-8 backdrop-blur-sm border border-white"
          >
            <div class="flex items-center gap-3">
              <PetalMark :size="32" :color="s.color" />
              <span class="text-xs font-mono text-stone-400">0{{ i + 1 }}</span>
            </div>
            <h3 class="mt-5 text-xl font-semibold text-stone-800">{{ s.title }}</h3>
            <p class="mt-2 text-stone-600 leading-relaxed">{{ s.body }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- quotes -->
    <section class="px-6 py-20">
      <div class="max-w-5xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-semibold text-stone-800 text-center">
          from people who keep one.
        </h2>
        <div class="mt-14 grid md:grid-cols-3 gap-6">
          <figure
            v-for="q in quotes"
            :key="q.name"
            class="bg-white/40 rounded-3xl p-8 border border-white/80"
          >
            <blockquote class="text-stone-700 italic leading-relaxed">
              "{{ q.text }}"
            </blockquote>
            <figcaption class="mt-6 text-sm text-stone-500">
              — {{ q.name }}, <span class="text-stone-400">{{ q.role }}</span>
            </figcaption>
          </figure>
        </div>
      </div>
    </section>

    <!-- gallery -->
    <section class="px-6 py-20">
      <div class="max-w-6xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-semibold text-stone-800 text-center">
          months that became flowers.
        </h2>
        <p class="mt-3 text-stone-600 text-center max-w-xl mx-auto">
          a handful of finished months. the legend on the side tells you what each color means.
        </p>
        <div class="mt-14 grid md:grid-cols-[14rem_1fr] gap-8 md:gap-10 items-start">
          <aside class="bg-white/60 rounded-3xl p-6 border border-white">
            <h3 class="text-xs font-mono uppercase tracking-[0.2em] text-stone-500">
              legend
            </h3>
            <ul class="mt-5 space-y-4">
              <li v-for="m in MOODS" :key="m.id" class="flex items-start gap-3">
                <PetalMark :size="28" :color="m.color" class="shrink-0 mt-0.5" />
                <div>
                  <div class="text-sm font-semibold text-stone-800">{{ m.label }}</div>
                  <div class="text-xs text-stone-500 leading-snug">{{ m.blurb }}</div>
                </div>
              </li>
            </ul>
          </aside>

          <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            <article
              v-for="g in gallery"
              :key="g.title"
              class="bg-white/60 rounded-3xl p-6 border border-white overflow-hidden"
            >
              <div class="aspect-square rounded-2xl bg-gradient-to-br from-white to-stone-50 ring-1 ring-stone-100 relative overflow-hidden">
                <ClientOnly>
                  <Flower3D
                    preview
                    :days-in-month="g.days"
                    :entries-by-day="g.entries"
                    :today-day="null"
                    :month-key="`gallery-${g.title}`"
                    :center-color="g.center"
                    class="w-full h-full"
                  />
                  <template #fallback>
                    <div class="w-full h-full flex items-center justify-center">
                      <FlowerMark :size="120" />
                    </div>
                  </template>
                </ClientOnly>
              </div>
              <h3 class="mt-5 text-lg font-semibold text-stone-800">{{ g.title }}</h3>
              <p class="text-sm text-stone-500">{{ g.sub }}</p>
            </article>
          </div>
        </div>
      </div>
    </section>

    <!-- manifesto -->
    <section class="px-6 py-24">
      <div class="max-w-2xl mx-auto text-center">
        <PetalMark :size="40" color="#FF9EBB" class="inline-block" />
        <h2 class="mt-6 text-3xl md:text-4xl font-semibold text-stone-800">
          why we built this.
        </h2>
        <div class="mt-8 space-y-5 text-lg text-stone-700 leading-relaxed">
          <p>
            most apps that ask how you feel want something from you.
            a streak. an upgrade. an opinion you can be sold against.
          </p>
          <p>
            petal does not. it is a small, slow place to put one feeling a day.
            no charts. no scores. no notifications dressed up as encouragement.
            after a month, you have a flower. that is all, and it turns out
            that is plenty.
          </p>
          <p class="text-stone-500 italic">
            — flo, who made this in his kitchen.
          </p>
        </div>
      </div>
    </section>

    <!-- pricing teaser -->
    <section class="px-6 py-20">
      <div class="max-w-3xl mx-auto bg-white/70 rounded-[2rem] p-10 md:p-14 border border-white text-center">
        <h2 class="text-3xl md:text-4xl font-semibold text-stone-800">
          free to plant.
        </h2>
        <p class="mt-4 text-stone-600 max-w-md mx-auto">
          keep the last 30 days for free, forever. three dollars a month if you want
          to keep every flower you ever grow.
        </p>
        <div class="mt-10">
          <NuxtLink
            to="/signup"
            data-petals class="inline-flex items-center gap-2 px-7 py-3.5 rounded-full bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition"
          >
            begin a flower
          </NuxtLink>
          <div class="mt-4">
            <NuxtLink to="/pricing" class="text-sm text-stone-500 hover:text-stone-800 underline-offset-4 hover:underline">
              see full pricing
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <Footer />
  </main>
</template>

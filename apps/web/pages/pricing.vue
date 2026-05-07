<script setup lang="ts">
useSeoMeta({
  title: 'pricing — petal',
  description:
    'free to plant. three dollars a month to keep every flower forever. no ads, no tracking, no funny business.',
})

const free = [
  'one petal a day',
  'all five moods',
  'last 30 days of history',
  'export the current month',
  'works on phone, tablet, laptop',
]

const bloom = [
  'everything in free',
  'unlimited history — every month, every petal',
  'export any month, any year',
  'private link to share a finished flower',
  'first to see new things we make',
]

const faqs = [
  {
    q: 'can i cancel anytime?',
    a: 'yes. one click in settings. your flowers stay, you just go back to the free tier.',
  },
  {
    q: 'where is my data stored?',
    a: 'hosted in europe, end-to-end encrypted.',
  },
  {
    q: 'can i export my entries?',
    a: 'yes. one button gives you a tidy json file with everything you have written. it is yours.',
  },
  {
    q: 'is there a phone app?',
    a: 'not yet. the website works well on phones. a real app may come if enough people ask for it.',
  },
  {
    q: 'do you train AI on my entries?',
    a: 'no, never. i mean — don\'t flatter yourself.',
  },
]
</script>

<template>
  <main class="px-6 py-20 md:py-28">
    <div class="max-w-4xl mx-auto text-center">
      <PetalMark :size="40" color="#FF9EBB" class="inline-block" />
      <h1 class="mt-6 text-4xl md:text-5xl font-semibold text-stone-800 tracking-tight">
        simple, kind pricing.
      </h1>
      <p class="mt-4 text-lg text-stone-600 max-w-xl mx-auto">
        try it free for as long as you want. pay a small amount only if you want every
        flower you grow to stay yours forever.
      </p>
    </div>

    <div class="mt-16 max-w-4xl mx-auto grid md:grid-cols-2 gap-6">
      <article class="bg-white/60 rounded-3xl p-8 border border-white">
        <div class="flex items-center gap-3">
          <PetalMark :size="28" color="#FFD93D" />
          <h2 class="text-xl font-semibold text-stone-800">free</h2>
        </div>
        <p class="mt-4 text-4xl font-semibold text-stone-800">$0</p>
        <p class="text-sm text-stone-500">forever, no card.</p>
        <ul class="mt-6 space-y-3 text-stone-700">
          <li v-for="f in free" :key="f" class="flex items-start gap-3">
            <span class="mt-2 w-1.5 h-1.5 rounded-full bg-stone-400 shrink-0" />
            <span>{{ f }}</span>
          </li>
        </ul>
      </article>

      <article class="bg-white rounded-3xl p-8 ring-2 ring-pink-200 relative">
        <span
          class="absolute -top-3 right-6 px-3 py-1 rounded-full bg-pink-200 text-pink-800 text-xs font-mono"
        >
          for keepers
        </span>
        <div class="flex items-center gap-3">
          <FlowerMark :size="32" />
          <h2 class="text-xl font-semibold text-stone-800">bloom</h2>
        </div>
        <p class="mt-4 text-4xl font-semibold text-stone-800">
          $3<span class="text-lg text-stone-500"> / month</span>
        </p>
        <p class="text-sm text-stone-500">or $30 a year. cancel anytime.</p>
        <ul class="mt-6 space-y-3 text-stone-700">
          <li v-for="f in bloom" :key="f" class="flex items-start gap-3">
            <span class="mt-2 w-1.5 h-1.5 rounded-full bg-pink-400 shrink-0" />
            <span>{{ f }}</span>
          </li>
        </ul>
      </article>
    </div>

    <section class="mt-24 max-w-2xl mx-auto">
      <h2 class="text-2xl md:text-3xl font-semibold text-stone-800 text-center">
        small questions, plain answers.
      </h2>
      <dl class="mt-10 divide-y divide-stone-200/60">
        <div v-for="f in faqs" :key="f.q" class="py-6">
          <dt class="font-semibold text-stone-800">{{ f.q }}</dt>
          <dd class="mt-2 text-stone-600 leading-relaxed">{{ f.a }}</dd>
        </div>
      </dl>
    </section>

    <section class="mt-20 max-w-2xl mx-auto text-center">
      <NuxtLink
        to="/signup"
        data-petals class="inline-flex items-center gap-2 px-7 py-3.5 rounded-full bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition"
      >
        plant the first petal
        <span aria-hidden="true">→</span>
      </NuxtLink>
    </section>

    <Footer />
  </main>
</template>

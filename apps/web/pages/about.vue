<script setup lang="ts">
useSeoMeta({
  title: 'about — petal',
  description:
    'petal is a small, slow journal made by one person who wanted somewhere quiet to put one feeling a day.',
})
</script>

<template>
  <main class="px-6 py-20 md:py-28">
    <article class="max-w-2xl mx-auto">
      <div class="text-center">
        <PetalMark :size="48" color="#FF9EBB" class="inline-block" />
        <h1 class="mt-6 text-4xl md:text-5xl font-semibold text-stone-800 tracking-tight">
          a small thing,<br />made slowly.
        </h1>
      </div>

      <div class="mt-14 space-y-6 text-lg text-stone-700 leading-relaxed">
        <p>
          one tuesday in march i tried to remember how the previous week went,
          and i could not. not the days, exactly — just the shape of them.
          whether it had been heavy or light. whether i had felt held or alone.
          i had used three different apps that year that promised to help me
          notice my own life, and they had each, in their own way, made me
          notice them instead.
        </p>

        <p>
          so i started drawing a flower. one petal a day, in colored pencil,
          on the inside cover of a paperback. five colors, picked once, never
          tweaked. by the end of april i had a small messy bloom that told me
          more about the month than any chart ever had. it lived where i could
          see it. it asked for nothing. when i looked at it, i remembered.
        </p>

        <p>
          petal is that, made into something other people can also use.
          the water lily was hand-modeled from a flower i photographed in
          a pond outside kyoto in 2024. the colors come from a paint set
          i still have. the writing area is one line, on purpose — the
          point is not to journal harder, it is to journal at all.
        </p>

        <p>
          there is no team. there is no roadmap. there is no investor on
          a slack telling me to gamify retention. there is one person, a
          quiet evening hobby, and a wish that you find this useful for
          a while. if you do, the three dollars a month help me keep the
          server running and the lights low.
        </p>

        <p class="text-stone-500 italic">
          — flo<br />
          <span class="text-sm">munich, by the window</span>
        </p>
      </div>

      <div class="mt-20 text-center">
        <NuxtLink
          to="/signup"
          data-petals class="inline-flex items-center gap-2 px-7 py-3.5 rounded-full bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition"
        >
          try it for a week
        </NuxtLink>
      </div>
    </article>

    <Footer />
  </main>
</template>

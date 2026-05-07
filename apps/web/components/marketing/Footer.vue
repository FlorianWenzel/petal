<template>
  <footer class="px-6 pt-20 pb-10">
    <div class="max-w-5xl mx-auto flex flex-col items-center gap-4">
      <FlowerMark :size="32" />
      <nav class="flex items-center gap-6 text-sm text-stone-500">
        <NuxtLink to="/about" class="hover:text-stone-800 transition">about</NuxtLink>
        <NuxtLink to="/pricing" class="hover:text-stone-800 transition">pricing</NuxtLink>
        <NuxtLink to="/privacy" class="hover:text-stone-800 transition">privacy</NuxtLink>
      </nav>
      <p class="text-xs text-stone-400 mt-2">made slowly, in munich · © petal</p>
    </div>
  </footer>
</template>

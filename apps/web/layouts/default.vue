<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

const { user, logout } = useAuth()
const route = useRoute()
const showNav = computed(() => !route.path.startsWith('/app'))
</script>

<template>
  <div class="min-h-full">
    <header v-if="showNav" class="px-6 py-5">
      <nav class="max-w-5xl mx-auto flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-2 text-xl font-semibold text-stone-800">
          <FlowerMark :size="28" />
          <span>petal</span>
        </NuxtLink>
        <div class="flex items-center gap-4 text-sm">
          <NuxtLink to="/pricing" class="text-stone-600 hover:text-stone-800 transition">pricing</NuxtLink>
          <NuxtLink to="/about" class="text-stone-600 hover:text-stone-800 transition">about</NuxtLink>
          <template v-if="user">
            <NuxtLink
              to="/app"
              data-petals class="px-4 py-2 rounded-full bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition"
            >
              my flower
            </NuxtLink>
            <button class="text-stone-600 hover:text-stone-800" @click="logout">log out</button>
          </template>
          <template v-else>
            <NuxtLink to="/login" class="text-stone-600 hover:text-stone-800">log in</NuxtLink>
            <NuxtLink
              to="/signup"
              data-petals class="px-4 py-2 rounded-full bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition"
            >
              sign up
            </NuxtLink>
          </template>
        </div>
      </nav>
    </header>
    <slot />
  </div>
</template>

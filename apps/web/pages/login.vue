<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

definePageMeta({ layout: 'auth' })
useSeoMeta({ title: 'log in — petal', robots: 'noindex' })

const { login } = useAuth()
const email = ref('')
const password = ref('')
const error = ref<string | null>(null)
const pending = ref(false)

async function submit() {
  error.value = null
  pending.value = true
  try {
    await login(email.value, password.value)
    await navigateTo('/app')
  } catch (e) {
    error.value = (e as Error).message ?? 'login failed'
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="max-w-md mx-auto px-6 py-24">
    <h1 class="text-3xl font-semibold text-stone-800 text-center">welcome back</h1>
    <form class="mt-10 space-y-4" @submit.prevent="submit">
      <input
        v-model="email"
        type="email"
        required
        placeholder="email"
        class="w-full px-4 py-3 rounded-2xl bg-white/70 border border-white outline-none focus:bg-white"
      />
      <input
        v-model="password"
        type="password"
        required
        placeholder="password"
        class="w-full px-4 py-3 rounded-2xl bg-white/70 border border-white outline-none focus:bg-white"
      />
      <button
        :disabled="pending"
        class="w-full px-4 py-3 rounded-2xl bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition disabled:opacity-50"
      >
        {{ pending ? '…' : 'log in' }}
      </button>
      <p v-if="error" class="text-sm text-rose-600 text-center">{{ error }}</p>
      <p class="text-sm text-stone-500 text-center">
        no account?
        <NuxtLink to="/signup" class="text-stone-800 underline">sign up</NuxtLink>
      </p>
    </form>
  </main>
</template>

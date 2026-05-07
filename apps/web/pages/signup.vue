<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

definePageMeta({ layout: 'auth' })
useSeoMeta({ title: 'sign up — petal', robots: 'noindex' })

const { signup } = useAuth()
const email = ref('')
const password = ref('')
const error = ref<string | null>(null)
const pending = ref(false)

async function submit() {
  error.value = null
  pending.value = true
  try {
    await signup(email.value, password.value)
    await navigateTo('/app')
  } catch (e) {
    error.value = (e as Error).message ?? 'signup failed'
  } finally {
    pending.value = false
  }
}
</script>

<template>
  <main class="max-w-md mx-auto px-6 py-24">
    <h1 class="text-3xl font-semibold text-stone-800 text-center">start journaling</h1>
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
        minlength="8"
        placeholder="password (min 8 chars)"
        class="w-full px-4 py-3 rounded-2xl bg-white/70 border border-white outline-none focus:bg-white"
      />
      <button
        :disabled="pending"
        class="w-full px-4 py-3 rounded-2xl bg-rose-brand text-rose-brand-ink hover:bg-rose-brand-hover transition disabled:opacity-50"
      >
        {{ pending ? '…' : 'create account' }}
      </button>
      <p v-if="error" class="text-sm text-rose-600 text-center">{{ error }}</p>
      <p class="text-sm text-stone-500 text-center">
        have an account?
        <NuxtLink to="/login" class="text-stone-800 underline">log in</NuxtLink>
      </p>
    </form>
  </main>
</template>

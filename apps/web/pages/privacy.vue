<script setup lang="ts">
useSeoMeta({
  title: 'privacy — petal',
  description:
    'what petal stores, what it does not, where it lives. plain english, no legalese.',
})
</script>

<template>
  <main class="px-6 py-20 md:py-28">
    <article class="max-w-2xl mx-auto">
      <div class="text-center">
        <PetalMark :size="40" color="#B5C7F7" class="inline-block" />
        <h1 class="mt-6 text-4xl md:text-5xl font-semibold text-stone-800 tracking-tight">
          your entries are yours.
        </h1>
        <p class="mt-4 text-stone-600">
          plain english, the way it should be.
        </p>
      </div>

      <section class="mt-14 space-y-10 text-stone-700 leading-relaxed">
        <div>
          <h2 class="text-xl font-semibold text-stone-800">what petal stores</h2>
          <ul class="mt-4 space-y-2">
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-stone-400 shrink-0" />
              your email, so you can sign in.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-stone-400 shrink-0" />
              a salted argon2 hash of your password. we never see the password.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-stone-400 shrink-0" />
              your mood and the line you write each day, encrypted before they ever leave your device. we hold the box, you hold the key.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-stone-400 shrink-0" />
              a session cookie so you do not have to log in every time.
            </li>
          </ul>
        </div>

        <div>
          <h2 class="text-xl font-semibold text-stone-800">what petal does not store</h2>
          <ul class="mt-4 space-y-2">
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-pink-300 shrink-0" />
              analytics. no google, no plausible, no mixpanel, no fathom.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-pink-300 shrink-0" />
              third-party trackers. there are zero scripts on this site that are not ours.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-pink-300 shrink-0" />
              your ip address, beyond the few seconds it takes to serve a request.
            </li>
            <li class="flex items-start gap-3">
              <span class="mt-2 w-1.5 h-1.5 rounded-full bg-pink-300 shrink-0" />
              ad ids, device fingerprints, or anything that follows you elsewhere.
            </li>
          </ul>
        </div>

        <div>
          <h2 class="text-xl font-semibold text-stone-800">where it lives</h2>
          <p class="mt-3">
            hosted in europe. end-to-end encrypted: the contents of your entries
            are scrambled on your device before being sent, and only your password
            can unscramble them. our server can store them, back them up, and serve
            them back — but it cannot read them. neither can we.
          </p>
          <p class="mt-3 text-sm text-stone-500">
            one practical consequence: if you forget your password, we cannot
            recover your old entries. that is the price of a real lock.
          </p>
        </div>

        <div>
          <h2 class="text-xl font-semibold text-stone-800">we never train AI on your entries</h2>
          <p class="mt-3">
            we couldn't if we wanted to — your entries are encrypted on the way out
            of your device. we also wouldn't if we could.
          </p>
        </div>

        <div>
          <h2 class="text-xl font-semibold text-stone-800">leaving</h2>
          <p class="mt-3">
            you can delete your account from settings. it is one button. it
            removes every entry, every session, your email, and the password
            hash. backups containing it expire within 30 days. nothing remains
            after that.
          </p>
        </div>

        <div>
          <h2 class="text-xl font-semibold text-stone-800">questions</h2>
          <p class="mt-3">
            write to <a href="mailto:hello@petal.app" class="text-stone-900 underline underline-offset-4">hello@petal.app</a>.
            a real person reads it, usually within a day or two.
          </p>
        </div>
      </section>

      <p class="mt-16 text-sm text-stone-400 text-center">
        last updated: 03 may 2026
      </p>
    </article>

    <Footer />
  </main>
</template>

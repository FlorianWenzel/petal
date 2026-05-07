<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

interface Particle {
  id: number
  x: number
  y: number
  vx: number
  vy: number
  rot: number
  vrot: number
  scale: number
  hue: number
  born: number
}

const particles = ref<Particle[]>([])
const MAX = 60
const LIFE_MS = 2600
const SPAWN_MIN_DIST_PX = 14
const HOVER_SPAWN_INTERVAL_MS = 110

let nextId = 0
let lastX = 0
let lastY = 0
let lastSpawn = 0
let raf = 0
let mounted = false
let hovering = false

function spawn(x: number, y: number) {
  if (particles.value.length >= MAX) {
    particles.value.shift()
  }
  const angle = Math.random() * Math.PI * 2
  const speed = 20 + Math.random() * 35
  particles.value.push({
    id: nextId++,
    x: x + (Math.random() - 0.5) * 8,
    y: y + (Math.random() - 0.5) * 8,
    vx: Math.cos(angle) * speed * 0.5,
    vy: Math.sin(angle) * speed * 0.5 - 18,
    rot: Math.random() * 360,
    vrot: (Math.random() - 0.5) * 180,
    scale: 0.55 + Math.random() * 0.55,
    hue: 330 + Math.random() * 30,
    born: performance.now(),
  })
}

const INTERACTIVE_SELECTOR = 'button, [role="button"], input[type="button"], input[type="submit"], [data-petals]'

function isOverInteractive(target: EventTarget | null): boolean {
  if (!(target instanceof Element)) return false
  return !!target.closest(INTERACTIVE_SELECTOR)
}

function onMove(e: PointerEvent) {
  const x = e.clientX
  const y = e.clientY
  const dx = x - lastX
  const dy = y - lastY
  lastX = x
  lastY = y
  hovering = isOverInteractive(e.target)
  if (!hovering) return
  if (Math.hypot(dx, dy) < SPAWN_MIN_DIST_PX) return
  const now = performance.now()
  if (now - lastSpawn < 30) return
  lastSpawn = now
  spawn(x, y)
}

function onLeave() {
  hovering = false
}

function tick() {
  if (!mounted) return
  const now = performance.now()
  const dt = 1 / 60

  if (hovering && now - lastSpawn > HOVER_SPAWN_INTERVAL_MS) {
    lastSpawn = now
    spawn(lastX, lastY)
  }

  particles.value = particles.value.filter((p) => {
    const age = now - p.born
    if (age > LIFE_MS) return false
    p.x += p.vx * dt
    p.y += p.vy * dt
    p.vy += 35 * dt
    p.vx *= 0.99
    p.rot += p.vrot * dt
    return true
  })
  raf = requestAnimationFrame(tick)
}

function lifeProgress(p: Particle): number {
  return Math.min(1, (performance.now() - p.born) / LIFE_MS)
}

onMounted(() => {
  mounted = true
  window.addEventListener('pointermove', onMove, { passive: true })
  window.addEventListener('pointerleave', onLeave, { passive: true })
  window.addEventListener('blur', onLeave)
  raf = requestAnimationFrame(tick)
})

onUnmounted(() => {
  mounted = false
  window.removeEventListener('pointermove', onMove)
  window.removeEventListener('pointerleave', onLeave)
  window.removeEventListener('blur', onLeave)
  cancelAnimationFrame(raf)
})
</script>

<template>
  <div class="petal-cursor" aria-hidden="true">
    <svg
      v-for="p in particles"
      :key="p.id"
      class="petal"
      :style="{
        transform: `translate3d(${p.x}px, ${p.y}px, 0) rotate(${p.rot}deg) scale(${p.scale})`,
        opacity: 1 - lifeProgress(p),
      }"
      width="18"
      height="18"
      viewBox="0 0 32 32"
    >
      <path
        d="M16 2 C 22 8, 24 16, 16 30 C 8 16, 10 8, 16 2 Z"
        :fill="`hsl(${p.hue}, 78%, 72%)`"
      />
      <path
        d="M16 6 C 19 11, 20 18, 16 26 C 12 18, 13 11, 16 6 Z"
        fill="white"
        opacity="0.4"
      />
    </svg>
  </div>
</template>

<style scoped>
.petal-cursor {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 9999;
}
.petal {
  position: absolute;
  top: 0;
  left: 0;
  margin-left: -9px;
  margin-top: -9px;
  will-change: transform, opacity;
}
</style>

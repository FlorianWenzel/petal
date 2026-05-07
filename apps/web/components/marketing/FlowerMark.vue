<script setup lang="ts">
withDefaults(defineProps<{ size?: number }>(), { size: 28 })

const petals = [
  { c: '#FFD93D', r: -90 },
  { c: '#FF9EBB', r: -30 },
  { c: '#FFA372', r: 30 },
  { c: '#B5C7F7', r: 90 },
  { c: '#4D6F9E', r: 150 },
  { c: '#FFC8DD', r: 210 },
]
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 64 64" aria-hidden="true">
    <g v-for="(p, i) in petals" :key="i" :transform="`rotate(${p.r} 32 32)`">
      <ellipse cx="32" cy="18" rx="6" ry="11" :fill="p.c" opacity="0.9" />
    </g>
    <circle cx="32" cy="32" r="6" fill="#FFF1B8" />
    <circle cx="32" cy="32" r="3" fill="#E8B86A" opacity="0.7" />
  </svg>
</template>

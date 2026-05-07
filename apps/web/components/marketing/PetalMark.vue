<script setup lang="ts">
withDefaults(defineProps<{ size?: number; color?: string }>(), {
  size: 24,
  color: '#FF9EBB',
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 32 32" aria-hidden="true">
    <path
      d="M16 2 C 22 8, 24 16, 16 30 C 8 16, 10 8, 16 2 Z"
      :fill="color"
      opacity="0.85"
    />
    <path
      d="M16 6 C 19 11, 20 18, 16 26 C 12 18, 13 11, 16 6 Z"
      fill="white"
      opacity="0.35"
    />
  </svg>
</template>

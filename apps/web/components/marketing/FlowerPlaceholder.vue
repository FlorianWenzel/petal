<script setup lang="ts">
withDefaults(defineProps<{ label?: string; size?: number }>(), {
  label: '3D flower',
  size: 480,
})
</script>

<template>
  <div class="relative flex items-center justify-center" :style="{ width: `${size}px`, height: `${size}px` }">
    <div
      class="absolute bottom-4 left-1/2 -translate-x-1/2 w-3/4 h-6 rounded-full bg-stone-900/10 blur-2xl"
    />
    <div
      class="relative w-full h-full rounded-full bg-gradient-to-b from-white/70 to-white/30 ring-1 ring-white/80 backdrop-blur-sm flex items-center justify-center overflow-hidden"
    >
      <FlowerMark :size="Math.floor(size * 0.55)" />
      <span
        class="absolute bottom-6 left-1/2 -translate-x-1/2 text-[10px] uppercase tracking-[0.2em] text-stone-500/60 font-mono"
      >
        {{ label }}
      </span>
    </div>
  </div>
</template>

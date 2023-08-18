<script lang="ts" setup>
defineSlots<{
  header(): any
  default(): any
}>()
</script>

<template>
  <div class="bg-Surface rounded">
    <div v-if="!!$slots.header"
      class="border-Overlay bg-Surface-100 overflow-hidden rounded-t border px-2 py-1 font-bold">
      <slot name="header"></slot>
    </div>
    <div
      class="border-x-Overlay first:border-t-Overlay last:border-b-Overlay overflow-auto border-x first:rounded-t first:border-t last:rounded-b last:border-b">
      <slot />
    </div>
  </div>
</template>

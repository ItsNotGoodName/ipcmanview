<script lang="ts" setup>
const colorMode = useColorMode()

function toggleColorMode() {
  if (colorMode.preference == "system") {
    colorMode.preference = "dark"
  } else if (colorMode.preference == "dark") {
    colorMode.preference = "light"
  } else {
    colorMode.preference = "system"
  }
}
</script>

<template>
  <ClientOnly>
    <template #fallback>
      <div class="bg-Overlay h-7 w-7 animate-pulse rounded-full" />
    </template>
    <button title="Theme" class="hover:text-Mauve hover:fill-Mauve flex items-center justify-center"
      @click="toggleColorMode">
      <Icon v-show="colorMode.preference == 'dark'" name="ri:moon-line" class="h-7 w-7" />
      <Icon v-show="colorMode.preference == 'light'" name="ri:sun-line" class="h-7 w-7" />
      <Icon v-show="colorMode.preference == 'system'" name="ri:contrast-line" class="h-7 w-7" />
    </button>
  </ClientOnly>
</template>

<style scoped></style>

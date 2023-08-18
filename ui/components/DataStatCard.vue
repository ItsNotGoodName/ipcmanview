<script lang="ts" setup>
import { type RouteLocationRaw } from ".nuxt/vue-router";
import { type AsyncData } from "nuxt/dist/app/composables/asyncData";

defineProps<{ data: AsyncData<unknown, unknown>, label: string, icon: string, to?: RouteLocationRaw }>()
</script>

<template>
  <UiCard>
    <div class="flex items-center gap-4 p-4">
      <NuxtLink :to="to">
        <Icon :name="icon" class="h-6 w-6 flex-shrink-0" />
      </NuxtLink>
      <div>
        <Icon v-if="data.pending.value" class="h-8 w-8 animate-spin" name="ri:loader-3-fill" />
        <Icon v-else-if="data.error.value" class="text-Red h-8 w-8" name="ri:error-warning-fill" />
        <div v-else class="text-2xl font-bold">{{ data.data.value }}</div>
        <div>{{ label }}</div>
      </div>
    </div>
  </UiCard>
</template>

<style scoped></style>

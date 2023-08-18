<script lang="ts" setup>
import { type AsyncDataRequestStatus } from "nuxt/dist/app/composables/asyncData";

const props = defineProps<{
  status: AsyncDataRequestStatus
}>()
</script>

<template>
  <div v-if="props.status == 'success'"
    class="bg-Green text-Crust flex flex-row items-center gap-2 overflow-auto rounded-xl p-2">
    <Icon class="h-5 w-5 flex-shrink-0" name="ri:check-fill" />
    <slot />
  </div>
  <div v-else-if="props.status == 'error'"
    class="bg-Red text-Crust flex flex-row items-center gap-2 overflow-auto rounded-xl p-2">
    <Icon class="h-5 w-5 flex-shrink-0" name="ri:error-warning-fill" />
    <slot />
  </div>
  <div v-else-if="props.status == 'pending'"
    class="bg-Blue text-Crust flex flex-row items-center gap-2 overflow-auto rounded-xl p-2">
    <Icon class="h-5 w-5 flex-shrink-0 animate-spin" name="ri:loader-3-fill" />
    <slot />
  </div>
  <div v-else-if="props.status == 'idle'"
    class="bg-Surface flex flex-row items-center gap-2 overflow-auto rounded-xl p-2">
    <Icon class="h-5 w-5 flex-shrink-0" name="ri:information-fill" />
    <slot />
  </div>
</template>

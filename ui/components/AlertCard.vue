<script lang="ts" setup>
import { Alert } from '~/stores/alert';

const props = defineProps<{ alert: Alert }>()
const emit = defineEmits(['close'])

const timeout = ref(props.alert.timeout)
const id = ref<NodeJS.Timer>()
const width = computed(() => Math.round((timeout.value / props.alert.timeout) * 100))

function stop() {
  if (id.value) {
    clearInterval(id.value)
    id.value = undefined
  }
}

function start() {
  if (!props.alert.timeout) {
    return
  }

  if (!id.value) {
    id.value = setInterval(() => {
      timeout.value -= .25
      if (timeout.value < 0) {
        emit("close")
        stop()
      }
    }, 250);
  }
}

onMounted(() => {
  start()
})

onUnmounted(() => {
  stop()
})
</script>

<template>
  <div
    class="bg-Surface text-Text border-Overlay focus:border-Overlay-200 pointer-events-auto flex flex-col rounded border shadow-lg"
    @focusin="stop" @focusout="start">
    <div class="flex items-center gap-2 p-2 text-lg font-bold" :class="{ 'border-inherit border-b': alert.message }">
      <Icon v-if="alert.type == 'info'" name="ri:information-line" class="text-Blue h-6 w-6 flex-shrink-0" />
      <Icon v-else-if="alert.type == 'error'" name="ri:error-warning-fill" class="text-Red h-6 w-6 flex-shrink-0" />
      <Icon v-else-if="alert.type == 'success'" name="ri:check-fill" class="text-Green h-6 w-6 flex-shrink-0" />
      <div class="flex-1 truncate">{{ alert.title }}</div>
      <button class="hover:bg-Surface-200 flex flex-shrink-0 items-center rounded-full" :tabindex="$attrs.tabindex + '.5'"
        @click="emit('close', ($event.target as HTMLInputElement).value)">
        <Icon name="ri:close-line" class="text-Red h-6 w-6" />
      </button>
    </div>
    <div class="p-2" v-if="alert.message">
      {{ alert.message }}
    </div>
    <div class="rounded border-2 transition-all duration-300 ease-linear"
      :class="{ 'border-Blue': alert.type == 'info', 'border-Red': alert.type == 'error', 'border-Green': alert.type == 'success' }"
      :style="{ width: width + '%' }" />
  </div>
</template>

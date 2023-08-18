<script lang="ts" setup>
defineProps<{ id?: string, label?: string, message?: string, color?: "success" | "error", required?: boolean, modelValue?: string }>()
defineEmits(['update:modelValue'])
defineOptions({
  inheritAttrs: false
})
</script>

<template>
  <div>
    <label :for="id" :class="{ 'text-Red': color == 'error', 'text-Green': color == 'success' }"
      class="mb-2 block overflow-hidden text-ellipsis">
      {{ label }}
      <span v-if="color == 'error' && required">*</span>
    </label>
    <input :value="modelValue" @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)" :id="id"
      :class="{ 'text-Red border-Red': color == 'error', 'text-Green border-Green': color == 'success', 'border-Overlay': !color }"
      class="bg-Base placeholder-Subtext block w-full rounded border p-2.5 text-sm disabled:opacity-70" v-bind="$attrs" />
    <p :class="{ 'text-Red': color == 'error', 'text-Green': color == 'success' }" class="mt-2 text-sm">
      {{ message }}
    </p>
  </div>
</template>

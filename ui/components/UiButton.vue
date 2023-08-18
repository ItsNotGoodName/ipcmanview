<script lang="ts" setup>
const props = withDefaults(defineProps<{
  size?: "md" | "sm",
  color?: "primary" | "secondary" | "error" | "success" | "info" | "warning"
}>(), {
  size: "md",
  color: "primary"
})

function sizeClass(): string {
  if (props.size == "md") {
    return "p-2"
  }
  return "text-sm p-1"
}

function colorClass(): string {
  if (props.color == "primary") {
    return "bg-Mauve text-Crust hover:bg-Mauve-200 disabled:bg-Mauve border-Mauve-200"
  }
  if (props.color == "secondary") {
    return "bg-Subtext text-Crust hover:bg-Subtext-100 disabled:bg-Subtext border-Subtext-100"
  }
  if (props.color == "error") {
    return "bg-Red text-Crust hover:bg-Red-200 disabled:bg-Red border-Red-200"
  }
  if (props.color == "success") {
    return "bg-Green text-Crust hover:bg-Green-200 disabled:bg-Green border-Green-200"
  }
  if (props.color == "info") {
    return "bg-Blue text-Crust hover:bg-Blue-200 disabled:bg-Blue border-Blue-200"
  }
  return "bg-Yellow text-Crust hover:bg-Yellow-200 disabled:bg-Yellow border-Yellow-200"
}

const classStyle = computed(() => sizeClass() + " " + colorClass() + " overflow-hidden border rounded disabled:opacity-70")
</script>

<template>
  <button :class="classStyle">
    <slot />
  </button>
</template>

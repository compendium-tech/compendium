<template>
  <button :type="type" :disabled="disabled" :class="[
    ' font-bold rounded-lg transition-all transform', buttonClasses,
    { ' opacity-50 cursor-not-allowed': disabled }]" v-bind="$attrs" @click="handleClick">
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from "vue"

interface BaseButtonProps {
  type?: "button" | "submit" | "reset"
  variant?: "primary" | "secondary" | "danger" | "outline"
  size?: "sm" | "md" | "lg"
  disabled?: boolean
  hoverEffect?: "translate" | "scale" | "none"
}

const props = withDefaults(defineProps<BaseButtonProps>(), {
  type: "button",
  variant: "primary",
  size: "sm",
  disabled: false,
  hoverEffect: "translate",
})

const emit = defineEmits(["click"])

const buttonClasses = computed(() => {
  let classes: string[] = []

  switch (props.variant) {
    case "primary":
      classes.push("bg-primary-600 hover:bg-primary-400 text-white")
      break
    case "secondary":
      classes.push("bg-gray-200 hover:bg-gray-300 text-gray-800")
      break
    case "danger":
      classes.push("bg-red-600 hover:bg-red-700 text-white")
      break
    case "outline":
      classes.push("border border-gray-300 text-gray-700 hover:bg-gray-50")
      break
  }

  switch (props.size) {
    case "sm":
      classes.push("px-2 py-2 text-sm")
      break
    case "md":
      classes.push("px-4 py-3 text-base")
      break
    case "lg":
      classes.push("px-8 py-3 text-lg")
      break
  }

  if (props.hoverEffect === "scale" && !props.disabled) {
    classes.push("hover:scale-105")
  }

  if (props.hoverEffect === "translate" && !props.disabled) {
    classes.push("hover:-translate-y-1")
  }

  return classes
})

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit("click", event)
  }
}
</script>

<template>
  <button :type="type" :disabled="disabled" :class="[
    'font-bold py-3 px-8 rounded-xl text-lg transition-all transform',
    buttonClasses,
    { 'opacity-50 cursor-not-allowed': disabled }
  ]" v-bind="$attrs" @click="handleClick">
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface BaseButtonProps {
  type?: 'button' | 'submit' | 'reset';
  variant?: 'primary' | 'secondary' | 'danger' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  hoverEffect?: 'scale' | 'none';
}

const props = withDefaults(defineProps<BaseButtonProps>(), {
  type: 'button',
  variant: 'primary',
  size: 'md',
  disabled: false,
  hoverEffect: 'scale',
});

const emit = defineEmits(['click']);

const buttonClasses = computed(() => {
  let classes: string[] = [];

  // Variant classes
  switch (props.variant) {
    case 'primary':
      classes.push('bg-primary-600 hover:bg-primary-700 text-white');
      break;
    case 'secondary':
      classes.push('bg-gray-200 hover:bg-gray-300 text-gray-800');
      break;
    case 'danger':
      classes.push('bg-red-600 hover:bg-red-700 text-white');
      break;
    case 'outline':
      classes.push('border border-gray-300 text-gray-700 hover:bg-gray-50');
      break;
  }

  switch (props.size) {
    case 'sm':
      classes.push('px-4 py-2 text-sm');
      break;
    case 'md':
      classes.push('px-6 py-3 text-lg');
      break;
    case 'lg':
      classes.push('px-10 py-4 text-xl');
      break;
  }

  if (props.hoverEffect === 'scale') {
    classes.push('hover:scale-105');
  }

  return classes;
});

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event);
  }
};
</script>

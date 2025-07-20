<template>
  <div class="relative w-full">
    <input :type="currentInputType" :class="mergedClasses" v-bind="$attrs" @input="updateValue" />
    <button v-if="isPasswordType" type="button" @click="togglePasswordVisibility"
      class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5"
      aria-label="Toggle password visibility">
      <Icon :icon="isPasswordVisible ? 'heroicons-solid:eye' : 'heroicons-solid:eye-slash'"
        class="w-5 h-5 text-gray-400" />
    </button>
    <p v-if="error" class="mt-2 text-sm text-red-600 whitespace-pre-line">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, toRefs, useAttrs } from 'vue';
import { Icon } from '@iconify/vue';

interface Props {
  error?: string;
  type?: string;
}

const props = defineProps<Props>();
const attrs = useAttrs();

const { error, type } = toRefs(props);

const isPasswordVisible = ref(false);

const baseInputClasses = [
  'block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900',
  'outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400',
  'focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600 sm:text-sm/6',
];

const emit = defineEmits(['update:modelValue', 'input']);

const updateValue = (event: Event): void => {
  const target = event.target as HTMLInputElement;
  emit('update:modelValue', target.value);
  emit('input', event);
};

const isPasswordType = computed<boolean>(() => type?.value === 'password');

const currentInputType = computed<string>(() => {
  return isPasswordType.value && isPasswordVisible.value ? 'text' : type?.value || 'text';
});

const mergedClasses = computed<string>(() => {
  const incomingClasses = attrs.class || '';
  const errorClass = props.error ? 'outline-red-600 focus:outline-red-600' : '';
  const paddingRightClass = isPasswordType.value ? 'pr-10' : '';

  return [
    ...baseInputClasses,
    errorClass,
    incomingClasses,
    paddingRightClass,
  ].filter(Boolean).join(' ');
});

const togglePasswordVisibility = (): void => {
  isPasswordVisible.value = !isPasswordVisible.value;
};
</script>

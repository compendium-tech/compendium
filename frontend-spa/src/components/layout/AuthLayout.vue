<template>
  <div class="flex min-h-screen flex-col items-center justify-center px-6 py-12 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <RouterLink to="/">
        <img src="/icon.svg" alt="Your Company" class="mx-auto h-20 w-auto" />
      </RouterLink>
      <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
        {{ headerText }}
      </h2>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <slot />
    </div>

    <p v-if="formKind !== AuthFormKind.Other" class="mt-10 text-center text-sm/6 text-gray-500">
      <template v-if="formKind == AuthFormKind.SignUpForm">
        Already a member?
        <RouterLink to="/auth/signin" class="font-semibold text-primary-600 hover:text-primary-500">Sign in to your
          account</RouterLink>
      </template>
      <template v-else>
        Not a member?
        <RouterLink to="/auth/signin" class="font-semibold text-primary-600 hover:text-primary-500">Create your
          account</RouterLink>
      </template>
    </p>
  </div>
</template>

<script lang="ts">
export enum AuthFormKind {
  SignInForm,
  SignUpForm,
  Other
}
</script>

<script setup lang="ts">
import { RouterLink } from "vue-router"

interface Props {
  headerText: string,
  formKind: AuthFormKind
}

defineProps<Props>()
</script>

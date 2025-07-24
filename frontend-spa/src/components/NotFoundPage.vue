<template>
  <section class="bg-white dark:bg-gray-900 ">
    <div class="container flex items-center min-h-screen px-6 py-12 mx-auto">
      <div class="flex flex-col items-center max-w-sm mx-auto text-center">
        <p class="p-3 text-sm font-medium text-primary-600 rounded-full bg-primary-100 dark:bg-gray-800">
          <Icon icon="streamline-freehand:server-error-404-not-found" class="w-9 h-9" />
        </p>
        <h1 class="mt-3 text-2xl font-semibold text-gray-800 dark:text-white md:text-3xl">Page not found</h1>
        <p class="mt-4 text-gray-500 text-lg dark:text-gray-400">The page you are looking for doesn't exist.</p>

        <div class="flex items-center mt-6 gap-x-3 w-auto">
          <RouterLink to="/">
            <BaseButton size="md" class="flex items-center space-x-2">
              <Icon icon="game-icons:return-arrow" />
              <span>Take me home</span>
            </BaseButton>
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router"
import BaseButton from "./ui/BaseButton.vue"
import { Icon } from "@iconify/vue";
</script>

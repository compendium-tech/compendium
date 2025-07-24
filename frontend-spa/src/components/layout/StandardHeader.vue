<template>
  <div class="fixed w-full bg-white/90 backdrop-blur-sm py-4 px-6 z-50">
    <div class="max-w-6xl mx-auto flex justify-between items-center">
      <RouterLink to="/">
        <div class="flex items-center">
          <div class="w-13 h-13 rounded-full flex items-center justify-center">
            <img src="/icon.svg" />
          </div>
          <span class="ml-1 text-xl font-semibold">compendium</span>
          <span class="text-xl font-bold text-primary-600">.</span>
        </div>
      </RouterLink>
      <template v-if="authStore.isAuthenticated">
        <RouterLink to="/dashboard">
          <Icon class="w-10 h-10 text-primary-600 transition-all transform hover:scale-105"
            icon="tdesign:user-circle-filled" />
        </RouterLink>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from "@iconify/vue"
import { useAuthStore } from "../../stores/auth.ts"

const authStore = useAuthStore()
</script>

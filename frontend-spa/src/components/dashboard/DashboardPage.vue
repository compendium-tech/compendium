<template>
  <StandardLayout>
    <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-4xl w-full bg-white p-8 space-y-8 rounded-lg mt-16">
        <div class="text-center">
          <h1 class="text-4xl font-extrabold text-gray-900">Your Dashboard</h1>
          <p class="mt-2 text-lg text-gray-600">Manage your account and active sessions.</p>
        </div>

        <AccountInfo />

        <div class="mt-8">
          <ActiveSessions />
        </div>
      </div>
    </div>
  </StandardLayout>
</template>

<script setup lang="ts">
import StandardLayout from '../layout/StandardLayout.vue'
import AccountInfo from './AccountInfoSection.vue'
import ActiveSessions from './ActiveSessionsSection.vue'
</script>

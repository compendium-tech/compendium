<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-800 mb-4">Active Sessions</h2>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-6">
      <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-blue-600"></div>
      <p class="mt-2 text-md text-blue-600">Loading sessions...</p>
    </div>
    <BaseTransitioningText>
      <div v-if="sessionsError" class="text-center text-red-600">{{ sessionsError }}</div>
    </BaseTransitioningText>

    <ul class="divide-y divide-gray-200 rounded-lg border border-gray-200">
      <li v-for="session in sessions" :key="session.id"
        class="p-4 flex flex-col sm:flex-row justify-between items-start sm:items-center">
        <div class="flex-1">
          <div class="text-gray-900 flex space-x-4 items-center">
            <span>{{ session.name }}</span>
            <span v-if="session.isCurrent"
              class="text-sm py-2 px-3 font-semibold text-green-800 bg-green-100 rounded-lg flex space-x-2">
              <Icon icon="material-symbols-light:ar-on-you-outline" class="h-5 w-5" />
              <span>Current</span>
            </span>
            <BaseButton v-if="!session.isCurrent" @click="removeSession(session.id)" size="sm" variant="secondary"
              class="flex space-x-2" :is-loading="removingSessionId === session.id">
              <Icon icon="mdi:close-circle-outline" class="h-5 w-5" />
              <span>Remove</span>
            </BaseButton>
            <BaseButton v-else @click="logout" size="sm" variant="secondary" class="flex space-x-2">
              <Icon icon="mdi:logout" class="h-5 w-5" />
              <span>Logout</span>
            </BaseButton>
          </div>
          <p class="text-sm text-gray-600">IP address: {{ session.ipAddress }}</p>
          <p class="text-sm text-gray-600">Operating System: {{ session.os ? session.os : 'Unknown' }}</p>
          <p class="text-sm text-gray-600">Device: {{ session.device ? session.device : 'Unknown' }}</p>
          <p class="text-sm text-gray-600">Location: {{ session.location ? session.location : 'Unknown' }}</p>
          <p class="text-sm text-gray-600">Logged in: {{ formatSessionDate(session.createdAt) }}</p>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, Ref } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { sessionService, Session } from '../../api/auth'
import { Icon } from '@iconify/vue'
import BaseButton from '../ui/BaseButton.vue'
import BaseTransitioningText from '../ui/BaseTransitioningText.vue'

const sessions: Ref<Session[]> = ref([])
const isLoading = ref(true)
const sessionsError: Ref<string | null> = ref(null)
const removingSessionId = ref<string | null>(null)

const authStore = useAuthStore()

const formatSessionDate = (dateString: string): string => {
  try {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch (e) {
    return 'Invalid Date'
  }
}

const fetchSessions = async () => {
  isLoading.value = true
  sessionsError.value = null
  try {
    const data = await sessionService.getSessions()
    sessions.value = data
  } catch (err) {
    sessionsError.value = err.message || 'Failed to fetch sessions.'
  } finally {
    isLoading.value = false
  }
}

const removeSession = async (sessionId: string) => {
  removingSessionId.value = sessionId
  try {
    await sessionService.deleteSession(sessionId)
    await fetchSessions()
  } catch (err) {
    sessionsError.value = err.message || 'Failed to remove session.'
  } finally {
    removingSessionId.value = null
  }
}

const logout = async () => {
  try {
    await sessionService.logout()
    authStore.clearSession()
    location.href = '/auth/signin'
  } catch (err) {
    sessionsError.value = err.message || 'Failed to logout.'
  }
}

onMounted(() => {
  fetchSessions()
})
</script>

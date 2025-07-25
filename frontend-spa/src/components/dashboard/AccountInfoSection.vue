<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-800 mb-4">Account Information</h2>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-6">
      <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
      <p class="mt-2 text-md text-primary-600">Loading account info...</p>
    </div>
    <div v-if="error" class="text-center text-red-600">{{ error }}</div>

    <div v-if="user && !isLoading">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">Email Address</label>
          <p class="mt-1 text-lg text-gray-900">{{ user.email }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Account Created</label>
          <p class="mt-1 text-lg text-gray-900">{{ formattedCreationDate }}</p>
        </div>
        <div class="col-span-1 md:col-span-2">
          <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
          <div class="mt-1 flex items-center space-x-2" v-if="!isEditingName">
            <span class="text-lg">{{ user.name }}</span>
            <BaseButton variant="outline" class="px-2 flex items-center" @click="toggleEditName" size="sm"
              hover-effect="scale">
              <Icon icon="mdi:pencil" class="h-5 w-5 mr-2" />
              <span>Edit</span>
            </BaseButton>
          </div>
          <div class="mt-1 flex items-center space-x-2" v-else>
            <BaseInput v-model="editableName" type="text" id="name" :disabled="isSavingName" />
            <BaseButton class="px-4 flex items-center" @click="toggleEditName" size="sm" hover-effect="scale"
              :is-loading="isSavingName">
              <Icon icon="mdi:content-save" class="h-5 w-5 mr-2" />
              <span>Save</span>
            </BaseButton>
          </div>
          <BaseTransitioningText>
            <p v-if="nameUpdateMessage" :class="nameUpdateSuccess ? 'text-green-600' : 'text-red-600'"
              class="mt-2 text-sm">
              {{ nameUpdateMessage }}
            </p>
          </BaseTransitioningText>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, Ref } from 'vue'
import { userService, AccountDetails } from '../../api/user'
import { Icon } from '@iconify/vue'
import BaseInput from '../ui/BaseInput.vue'
import BaseButton from '../ui/BaseButton.vue'
import BaseTransitioningText from '../ui/BaseTransitioningText.vue'

const user: Ref<AccountDetails | null> = ref(null)
const isLoading = ref(true)
const error: Ref<string | null> = ref(null)
const isEditingName = ref(false)
const editableName = ref('')
const nameUpdateMessage = ref('')
const nameUpdateSuccess = ref(false)
const isSavingName = ref(false)

const formattedCreationDate = computed(() => {
  if (user.value && user.value.createdAt) {
    return new Date(user.value.createdAt).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }
  return 'N/A'
})

const fetchUserData = async () => {
  isLoading.value = true
  error.value = null
  try {
    const data = await userService.getAccountDetails()
    user.value = data
    editableName.value = data.name
  } catch (err: any) {
    error.value = err.message || 'Failed to fetch user data.'
  } finally {
    isLoading.value = false
  }
}

const toggleEditName = async () => {
  if (!user.value) return

  if (isEditingName.value) {
    isSavingName.value = true
    try {
      const updatedUser = await userService.updateName(editableName.value)
      user.value.name = updatedUser.name
      nameUpdateMessage.value = 'Name updated successfully!'
      nameUpdateSuccess.value = true
    } catch (err) {
      nameUpdateMessage.value = 'Failed to update name: ' + (err.message || 'Unknown error')
      nameUpdateSuccess.value = false
    } finally {
      isSavingName.value = false
      setTimeout(() => {
        nameUpdateMessage.value = ''
      }, 3000);
    }
  }
  isEditingName.value = !isEditingName.value
}

onMounted(() => {
  fetchUserData()
})
</script>

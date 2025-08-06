<template>
  <div>
    <h2 class="text-2xl font-semibold text-gray-800 mb-4">Account Information</h2>

    <div v-if="isLoadingAccount" class="flex flex-col items-center justify-center py-6">
      <div class="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-primary-600"></div>
      <p class="mt-2 text-md text-primary-600">Loading account info...</p>
    </div>
    <div v-if="globalError" class="text-center text-red-600">{{ globalError }}</div>

    <div v-if="user && !isLoadingAccount">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">Email Address</label>
          <p class="mt-1 text-lg text-gray-900">{{ user.email }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Creation Date</label>
          <p class="mt-1 text-lg text-gray-900">{{ dateToString(user.createdAt) }}</p>
        </div>
        <div class="col-span-1 md:col-span-2">
          <label class="block text-sm font-medium text-gray-700">Name</label>
          <div class="flex space-x-2">
            <BaseInput v-model="editableName" type="text" id="name" class="text-lg" :disabled="!isEditingName" />
            <BaseButton :variant="isEditingName ? 'primary' : 'outline'" class="space-x-2 flex items-center"
              @click="toggleEditName" hover-effect="scale" :is-loading="isSavingName">
              <Icon :icon="isEditingName ? 'mdi:content-save' : 'mdi:pencil'" class="h-5 w-5" />
              <span>{{ isEditingName ? 'Save' : 'Edit' }}</span>
            </BaseButton>
          </div>
          <BaseTransitioningText>
            <p v-if="nameUpdateError" :class="'text-red-600'" class="mt-2 text-sm">
              {{ nameUpdateError }}
            </p>
          </BaseTransitioningText>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, Ref } from 'vue'
import { userService, AccountDetails } from '../../api/user.ts'
import { Icon } from '@iconify/vue'
import BaseInput from '../ui/BaseInput.vue'
import BaseButton from '../ui/BaseButton.vue'
import BaseTransitioningText from '../ui/BaseTransitioningText.vue'
import { dateToString } from "../../utils/date"
import { ApiError } from '../../api/base.ts'

const user: Ref<AccountDetails | null> = ref(null)

const isLoadingAccount = ref(true)
const isEditingName = ref(false)
const isSavingName = ref(false)

const globalError = ref('')
const editableName = ref('')
const nameUpdateError = ref('')

const fetchUserData = async () => {
  isLoadingAccount.value = true
  globalError.value = ''

  try {
    const data = await userService.getAccountDetails()
    user.value = data
    editableName.value = data.name
  } catch (error) {
    if (error instanceof ApiError)
      globalError.value = error.message
  } finally {
    isLoadingAccount.value = false
  }
}

const toggleEditName = async () => {
  if (!user.value) return

  if (!isEditingName.value) {
    editableName.value = user.value.name
    isEditingName.value = true
    return
  }

  isSavingName.value = true

  try {
    const updatedUser = await userService.updateName(editableName.value)
    user.value.name = updatedUser.name

    isEditingName.value = false
  } catch (error) {
    if (error instanceof ApiError)
      nameUpdateError.value = error.message
  } finally {
    isSavingName.value = false

    setTimeout(() => nameUpdateError.value = '', 3000)
  }
}

onMounted(fetchUserData)
</script>

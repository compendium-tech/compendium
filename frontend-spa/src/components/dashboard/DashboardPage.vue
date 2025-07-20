<template>
  <StandardLayout>
    <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-4xl w-full bg-white p-8 space-y-8">
        <div class="text-center">
          <h1 class="text-4xl font-extrabold text-gray-900">Your Dashboard</h1>
          <p class="mt-2 text-lg text-gray-600">Manage your account and subscription settings.</p>
        </div>

        <div v-if="isLoading" class="flex flex-col items-center justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-t-4 border-b-4 border-primary-600"></div>
          <p class="mt-3 text-lg text-primary-600">Loading user data...</p>
        </div>
        <div v-if="error" class="text-center text-red-600">Error: {{ error }}</div>

        <div v-if="user && !isLoading">
          <div class="bg-gray-50 p-6 rounded-lg shadow-inner">
            <h2 class="text-2xl font-semibold text-gray-800 mb-4">Account Information</h2>
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
                <div class="mt-1 flex rounded-lg shadow-sm">
                  <input v-model="editableName" type="text" id="name" :disabled="!isEditingName"
                    class="flex-1 block w-full rounded-l-md border-primary-100 focus:border-primary-100 sm:text-lg p-2.5 disabled:bg-gray-100 disabled:text-gray-600" />
                  <button @click="toggleEditName"
                    class="flex items-center text-lg px-4 py-2 border border-transparent text-sm font-medium rounded-r-md shadow-sm text-white bg-primary-600 hover:bg-primary-400 transition-colors duration-200 ">
                    <Icon :icon="isEditingName ? 'mdi:content-save' : 'mdi:pencil'" class="h-5 w-5 mr-2" />
                    <span>{{ isEditingName ? 'Save' : 'Edit' }}</span>
                  </button>
                </div>
                <p v-if="nameUpdateMessage" :class="nameUpdateSuccess ? 'text-green-600' : 'text-red-600'"
                  class="mt-2 text-sm">{{ nameUpdateMessage }}</p>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 p-6 rounded-lg shadow-inner mt-8">
            <h2 class="text-2xl font-semibold text-gray-800 mb-4">Subscription Status</h2>
            <div class="flex flex-col md:flex-row items-center justify-between gap-4">
              <div class="text-lg text-gray-900">
                Status: <span :class="subscriptionStatusClass">{{ user.subscriptionStatus }}</span>
                <span v-if="user.subscriptionExpires"> (Expires: {{ formattedSubscriptionExpiry }})</span>
              </div>
              <button @click="handleSubscribe"
                class="flex items-center text-lg px-4 py-2 border border-transparent text-sm font-medium rounded-r-md shadow-sm text-white bg-primary-600 hover:bg-primary-400 transition-colors duration-200 ">
                <Icon icon="famicons:card-outline" class="w-6 h-6 mr-2" />
                <span>Enter billing information</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="showChangePasswordModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-75 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">Change Password</h3>
          <p class="text-gray-700 mb-6">This is a placeholder for the change password form. In a real application, you
            would have input fields for current and new passwords here.</p>
          <div class="flex justify-end space-x-3">
            <button @click="showChangePasswordModal = false"
              class="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-50">Cancel</button>
            <button @click="handleChangePassword"
              class="px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700">Confirm Change</button>
          </div>
        </div>
      </div>

      <div v-if="showRemoveAccountModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-75 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">Remove Account</h3>
          <p class="text-gray-700 mb-6">Are you sure you want to remove your account? This action cannot be undone.</p>
          <div class="flex justify-end space-x-3">
            <button @click="showRemoveAccountModal = false"
              class="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-50">Cancel</button>
            <button @click="handleRemoveAccount"
              class="px-4 py-2 rounded-lg bg-red-600 text-white hover:bg-red-700">Yes,
              Remove Account</button>
          </div>
        </div>
      </div>
    </div>
  </StandardLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { userService } from "../../api";
import StandardLayout from "../layout/StandardLayout.vue";
import { Icon } from '@iconify/vue';

const user = ref(null);
const isLoading = ref(true);
const error = ref(null);
const isEditingName = ref(false);
const editableName = ref('');
const nameUpdateMessage = ref('');
const nameUpdateSuccess = ref(false);

const showChangePasswordModal = ref(false);
const showRemoveAccountModal = ref(false);

const formattedCreationDate = computed(() => {
  if (user.value && user.value.createdAt) {
    return new Date(user.value.createdAt).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
  return 'N/A';
});

const formattedSubscriptionExpiry = computed(() => {
  if (user.value && user.value.subscriptionExpires) {
    return new Date(user.value.subscriptionExpires).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
  return 'N/A';
});

const subscriptionStatusClass = computed(() => {
  if (user.value) {
    switch (user.value.subscriptionStatus) {
      case 'active':
        return 'text-green-600 font-bold';
      case 'trialing':
        return 'text-primary-600 font-bold';
      case 'inactive':
        return 'text-red-600 font-bold';
      default:
        return 'text-gray-600';
    }
  }
  return '';
});

const fetchUserData = async () => {
  isLoading.value = true;
  error.value = null;
  const startTime = Date.now(); // Record the start time

  try {
    const data = (await userService.getAccountDetails()).data;
    user.value = data;
    editableName.value = data.name;
  } catch (err) {
    error.value = err.message || 'Failed to fetch user data.';
  } finally {
    const elapsedTime = Date.now() - startTime;
    const minimumLoadTime = 1000; // 1 second in milliseconds

    if (elapsedTime < minimumLoadTime) {
      setTimeout(() => {
        isLoading.value = false;
      }, minimumLoadTime - elapsedTime);
    } else {
      isLoading.value = false;
    }
  }
};

const toggleEditName = async () => {
  if (isEditingName.value) {
    try {
      await userService.updateName(editableName.value);
      user.value.name = editableName.value;

      nameUpdateMessage.value = 'Name updated successfully!';
      nameUpdateSuccess.value = true;
    } catch (err) {
      nameUpdateMessage.value = 'Failed to update name: ' + (err.message || 'Unknown error');
      nameUpdateSuccess.value = false;
    } finally {
      setTimeout(() => {
        nameUpdateMessage.value = '';
      }, 3000);
    }
  }
  isEditingName.value = !isEditingName.value;
};

const handleChangePassword = () => {
  console.log('Change Password button clicked. Implement password change form/logic here.');
  showChangePasswordModal.value = false;
};

const handleRemoveAccount = async () => {
  showRemoveAccountModal.value = false;
};

const handleSubscribe = async () => {
  // Logic for handling subscription goes here
};

onMounted(() => {
  fetchUserData();
});
</script>

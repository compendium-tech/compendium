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
        <div v-if="globalError" class="text-center text-red-600">{{ globalError }}</div>

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
                  <BaseButton class="flex items-center" @click="toggleEditName" variant="primary" size="sm"
                    :hover-effect="isEditingName ? 'none' : 'scale'">
                    <Icon :icon="isEditingName ? 'mdi:content-save' : 'mdi:pencil'" class="h-5 w-5 mr-2" />
                    <span>{{ isEditingName ? "Save" : "Edit" }}</span>
                  </BaseButton>
                </div>
                <p v-if="nameUpdateMessage" :class="nameUpdateSuccess ? ' text-green-600' : 'text-red-600'"
                  class=" mt-2 text-sm">{{ nameUpdateMessage }}</p>
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
              <BaseButton class="flex items-center" @click="handleSubscribe" variant="primary" size="sm">
                <Icon icon="famicons:card-outline" class="w-6 h-6 mr-2" />
                <span>Enter billing information</span>
              </BaseButton>
            </div>
          </div>
        </div>
      </div>
    </div>
  </StandardLayout>
</template>

<script setup>
import { ref, onMounted, computed } from "vue"
import { userService } from "../../api.ts"
import StandardLayout from "../layout/StandardLayout.vue"
import { Icon } from "@iconify/vue"
import BaseButton from "../ui/BaseButton.vue"

const user = ref(null)
const isLoading = ref(true)
const globalError = ref(null)
const isEditingName = ref(false)
const editableName = ref("")
const nameUpdateMessage = ref("")
const nameUpdateSuccess = ref(false)

const formattedCreationDate = computed(() => {
  if (user.value && user.value.createdAt) {
    return new Date(user.value.createdAt).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric"
    })
  }
  return "N/A"
})

const formattedSubscriptionExpiry = computed(() => {
  if (user.value && user.value.subscriptionExpires) {
    return new Date(user.value.subscriptionExpires).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric"
    })
  }
  return "N/A"
})

const subscriptionStatusClass = computed(() => {
  if (user.value) {
    switch (user.value.subscriptionStatus) {
      case "active":
        return "text-green-600 font-bold"
      case "trialing":
        return "text-primary-600 font-bold"
      case "inactive":
        return "text-red-600 font-bold"
      default:
        return "text-gray-600"
    }
  }
  return ""
})

const fetchUserData = async () => {
  isLoading.value = true
  globalError.value = null
  const startTime = Date.now()

  try {
    const data = await userService.getAccountDetails()
    user.value = data
    editableName.value = data.name
  } catch (err) {
    globalError.value = err.message || "Failed to fetch user data."
  } finally {
    const elapsedTime = Date.now() - startTime
    const minimumLoadTime = 1000

    if (elapsedTime < minimumLoadTime) {
      setTimeout(() => {
        isLoading.value = false
      }, minimumLoadTime - elapsedTime)
    } else {
      isLoading.value = false
    }
  }
}

const toggleEditName = async () => {
  if (isEditingName.value) {
    try {
      await userService.updateName(editableName.value)
      user.value.name = editableName.value

      nameUpdateMessage.value = "Name updated successfully!"
      nameUpdateSuccess.value = true
    } catch (err) {
      nameUpdateMessage.value = "Failed to update name: " + (err.message || "Unknown error")
      nameUpdateSuccess.value = false
    } finally {
      setTimeout(() => {
        nameUpdateMessage.value = ""
      }, 3000)
    }
  }
  isEditingName.value = !isEditingName.value
}

const handleSubscribe = () => {
  if (window.Paddle) {
    window.Paddle.Checkout.open({
      product: "pro_01k0az151b5zeh8a2974yvr0gx",
      customer: {
        email: user.value ? user.value.email : "",
      },
      items: [
        {
          priceId: "pri_01k0az35r3w64r2qmmzvc3rsgd",
          quantity: 1
        }
      ],
      successCallback: (data) => {
        console.log("Paddle checkout successful:", data)
      },
      closeCallback: () => {
        console.log("Paddle checkout closed.")
      }
    })
  } else {
    console.error("Paddle.js not loaded!")
    globalError.value = "Billing service is not available. Please try again later."
  }
}

onMounted(() => {
  if (window.Paddle) {
    window.Paddle.Environment.set("sandbox")
    window.Paddle.Setup({
      token: "test_bf5c18ea62fd1d30c00bc5c2821",
      debug: true
    })
  } else {
    console.error("Paddle.js script not found. Make sure it is included in your index.html")
  }
  fetchUserData()
})
</script>

<template>
  <AuthLayout :header-text="headerText" :form-kind="AuthFormKind.SignUpForm">
    <template v-if="state === State.Credentials">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="name" class="block text-sm/6 font-medium text-gray-900">Full Name</label>
          <div class="mt-2">
            <BaseInput id="name" type="text" v-model.trim="name" required autocomplete="name" placeholder="John Doe"
              @input="validateField('name')" />
            <BaseTransitioningText>
              <p v-if="validationErrors.name" class="mt-2 text-sm text-red-600">{{ validationErrors.name }}
              </p>
            </BaseTransitioningText>
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
          <div class="mt-2">
            <BaseInput id="email" type="email" v-model.trim="email" required autocomplete="email"
              placeholder="johndoe@gmail.com" input="validateField('email')" />
            <BaseTransitioningText>
              <p v-if="validationErrors.email" class="mt-2 text-sm text-red-600">{{ validationErrors.email }}
              </p>
            </BaseTransitioningText>
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
          <div class="mt-2">
            <BaseInput id="password" type="password" v-model="password" required autocomplete="new-password"
              placeholder="A strong password" @input="validateField('password')" />
            <BaseTransitioningText>
              <p v-if="validationErrors.password" class="mt-2 text-sm text-red-600 whitespace-pre-line">{{
                validationErrors.password }}</p>
            </BaseTransitioningText>
          </div>
        </div>

        <div>
          <BaseButton class="w-full" type="submit" :disabled="isLoading || !isCredentialsFormValid">
            Sign Up
            <span v-if="isLoading"
              class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent self-center"></span>
          </BaseButton>
        </div>
        <BaseTransitioningText>
          <p v-if="globalError" class="text-red-600 text-center text-sm mt-4">{{ globalError }}</p>
        </BaseTransitioningText>
      </form>
    </template>

    <template v-else-if="state === State.Mfa">
      <p class="mt-4 text-sm/6 text-gray-600">A One-Time Password has been sent to <span
          class="font-medium text-primary-600">{{ email }}</span>.</p>
      <form @submit.prevent="verifyMfa" class="space-y-6 mt-6">
        <div>
          <div class="mt-2">
            <BaseInput id="otp" type="text" v-model.trim="otp" placeholder="Enter 6-digit verification code" required
              maxlength="6" @input="validateField('otp')" />
            <BaseTransitioningText>
              <p v-if="validationErrors.otp" class="mt-2 text-sm text-red-600">{{ validationErrors.otp }}</p>
            </BaseTransitioningText>
          </div>
        </div>
        <div>
          <BaseButton class="w-full" type="submit" :disabled="isLoadingMfa || !isMfaFormValid">
            Verify Account
            <span v-if="isLoadingMfa"
              class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent self-center"></span>
          </BaseButton>
        </div>
        <BaseTransitioningText>
          <p v-if="globalError" class="text-red-600 text-center text-sm mt-4">{{ globalError }}</p>
        </BaseTransitioningText>
      </form>

      <div class="mt-6 space-x-3 flex">
        <BaseButton variant="secondary" @click="resendOtp" :disabled="countdown > 0 || isLoading">
          Resend Code <span v-if="countdown > 0">({{
            countdown }}s)</span>
        </BaseButton>
        <BaseButton variant="outline" @click="goBackToForm" class="font-semibold text-gray-600 hover:text-gray-500">Go
          Back</BaseButton>
      </div>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from "vue"
import { useRouter } from "vue-router"
import AuthLayout, { AuthFormKind } from "../layout/AuthLayout.vue"
import { authService } from "../../api/auth"
import { ApiError } from "../../api/base"
import { useAuthStore } from "../../stores/auth"
import { isEmailValid, isPasswordValid, isSixDigitCodeValid } from "../../utils/validation"
import BaseInput from "../ui/BaseInput.vue"
import BaseButton from "../ui/BaseButton.vue"
import BaseTransitioningText from "../ui/BaseTransitioningText.vue"

enum State {
  Credentials,
  Mfa,
}

const router = useRouter()
const authStore = useAuthStore()

const name = ref<string>("")
const email = ref<string>("")
const password = ref<string>("")
const otp = ref<string>("")

const isLoading = ref<boolean>(false)
const isLoadingMfa = ref<boolean>(false)
const globalError = ref<string>("")
const validationErrors = ref<Record<string, string>>({})
const state = ref<State>(State.Credentials)

const countdown = ref<number>(0)
let countdownTimer: number | undefined = undefined

const headerText = computed<string>(() => {
  switch (state.value) {
    case State.Mfa:
      return "Verify Your Email Address"
    case State.Credentials:
      return "Create New Account"
  }
})

const isCredentialsFormValid = computed<boolean>(() => {
  return name.value.length >= 1 && name.value.length <= 100 &&
    isEmailValid(email.value) &&
    isPasswordValid(password.value)
})

const isMfaFormValid = computed<boolean>(() => {
  return isSixDigitCodeValid(otp.value)
})

const validateField = (fieldName: string): void => {
  switch (fieldName) {
    case "name":
      if (!name.value.trim()) {
        validationErrors.value.name = "Full Name is required."
      } else if (name.value.length > 100) {
        validationErrors.value.name = "Full Name cannot exceed 100 characters."
      } else {
        delete validationErrors.value.name
      }
      break
    case "email":
      if (!email.value.trim()) {
        validationErrors.value.email = "Email address is required."
      } else if (!isEmailValid(email.value)) {
        validationErrors.value.email = "Please enter a valid email address."
      } else {
        delete validationErrors.value.email
      }
      break
    case "password":
      if (!password.value) {
        validationErrors.value.password = "Password is required."
      } else if (!isPasswordValid(password.value)) {
        validationErrors.value.password =
          "Password must be between 6 and 100 characters long, and contain at least one uppercase letter, one lowercase letter, one digit, and one special character (e.g., !@#$%^&*)."
      } else {
        delete validationErrors.value.password
      }
      break
    case "otp":
      if (!otp.value.trim()) {
        validationErrors.value.otp = "Verification code is required."
      } else if (!isSixDigitCodeValid(otp.value)) {
        validationErrors.value.otp = "Verification code must be a 6-digit number."
      } else {
        delete validationErrors.value.otp
      }
      break
  }
}

const validateCredentialsForm = (): boolean => {
  validationErrors.value = {}
  let isValid = true

  validateField("name")
  validateField("email")
  validateField("password")

  if (Object.keys(validationErrors.value).length > 0) {
    isValid = false
  }

  return isValid
}

const validateMfaForm = (): boolean => {
  validationErrors.value = {}
  let isValid = true
  validateField("otp")

  if (Object.keys(validationErrors.value).length > 0) {
    isValid = false
  }
  return isValid
}

const clearMessages = (): void => {
  globalError.value = ""
  validationErrors.value = {}
}

const startCountdown = (): void => {
  countdown.value = 60

  if (countdownTimer) {
    clearInterval(countdownTimer)
  }

  countdownTimer = setInterval(() => {
    if (countdown.value > 0) {
      countdown.value--
    } else {
      clearInterval(countdownTimer)
      countdownTimer = undefined
    }
  }, 1000)
}

const goBackToForm = (): void => {
  state.value = State.Credentials
  otp.value = ""

  clearMessages()

  if (countdownTimer) {
    clearInterval(countdownTimer)
  }

  countdown.value = 0
}

const handleSubmit = async (): Promise<void> => {
  clearMessages()
  if (!validateCredentialsForm()) {
    return
  }

  isLoading.value = true

  try {
    await authService.signUp(name.value, email.value, password.value)

    state.value = State.Mfa
    startCountdown()
  } catch (error) {
    if (error instanceof ApiError)
      globalError.value = error.message
  } finally {
    isLoading.value = false
  }
}

const verifyMfa = async (): Promise<void> => {
  clearMessages()
  if (!validateMfaForm()) {
    return
  }

  isLoadingMfa.value = true

  try {
    const response = await authService.verifyMfaSignUp(email.value, otp.value)

    authStore.setSession(response.accessTokenExpiresAt)
    router.push("/dashboard")
  } catch (error) {
    if (error instanceof ApiError)
      globalError.value = error.message
  } finally {
    isLoadingMfa.value = false
  }
}

const resendOtp = async (): Promise<void> => {
  clearMessages()
  if (countdown.value > 0 || isLoading.value) {
    return
  }

  isLoading.value = true

  try {
    await authService.signUp(name.value, email.value, password.value)
    startCountdown()
  } catch (error) {
    if (error instanceof ApiError)
      globalError.value = error.message
  } finally {
    isLoading.value = false
  }
}

onUnmounted(() => clearInterval(countdownTimer))
</script>

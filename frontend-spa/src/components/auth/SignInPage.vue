<template>
  <AuthLayout :header-text="headerText" :form-kind="AuthFormKind.SignInForm">
    <template v-if="state === State.Credentials">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
          <div class="mt-2">
            <input id="email" type="email" v-model.trim="email" required autocomplete="email"
              placeholder="johndoe@gmail.com" :class="inputClasses" @input="validateField('email')" />
            <p v-if="validationErrors.email" class="mt-2 text-sm text-red-600">{{ validationErrors.email }}</p>
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
          <div class="mt-2">
            <input id="password" type="password" v-model="password" required autocomplete="current-password"
              placeholder="A strong password" :class="inputClasses" @input="validateField('password')" />
            <p v-if="validationErrors.password" class="mt-2 text-sm text-red-600 whitespace-pre-line">{{
              validationErrors.password }}</p>
          </div>
        </div>

        <div>
          <button type="submit" :disabled="isLoading || !isCredentialsFormValid"
            class="flex w-full justify-center rounded-md bg-primary-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-primary-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:bg-primary-400 disabled:cursor-not-allowed">
            Sign In
            <span v-if="isLoading"
              class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent self-center"></span>
          </button>
        </div>
        <p v-if="serverError" class="text-red-600 text-center text-sm mt-4">{{ serverError }}</p>
      </form>
    </template>

    <template v-else-if="state === State.Mfa">
      <p class="mt-4 text-sm/6 text-gray-600">A One-Time Password has been sent to <span
          class="font-medium text-primary-600">{{ email }}</span>.</p>
      <form @submit.prevent="verifyMfa" class="space-y-6 mt-6">
        <div>
          <div class="mt-2">
            <input id="otp" type="text" v-model.trim="otp" placeholder="Enter 6-digit verification code" required
              maxlength="6" :class="inputClasses" @input="validateField('otp')" />
            <p v-if="validationErrors.otp" class="mt-2 text-sm text-red-600">{{ validationErrors.otp }}</p>
          </div>
        </div>
        <div>
          <button type="submit" :disabled="isLoadingMfa || !isMfaFormValid"
            class="flex w-full justify-center rounded-md bg-primary-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-primary-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:bg-primary-400 disabled:cursor-not-allowed">
            Verify Account
            <span v-if="isLoadingMfa"
              class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent self-center"></span>
          </button>
        </div>
        <p v-if="serverError" class="text-red-600 text-center text-sm mt-4">{{ serverError }}</p>
      </form>

      <div class="mt-6 text-center text-sm/6 space-y-3">
        <button @click="resendOtp" :disabled="countdown > 0 || isLoading"
          class="font-semibold text-primary-600 hover:text-primary-500 disabled:text-gray-400 disabled:cursor-not-allowed">
          Resend Code <span v-if="countdown > 0">({{ countdown }}s)</span>
        </button>
        <p>
          <button @click="goBackToForm" class="font-semibold text-gray-600 hover:text-gray-500">Go Back</button>
        </p>
      </div>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import AuthLayout, { AuthFormKind } from '../layout/AuthLayout.vue'
import { authService } from '../../api'
import { useAuthStore } from '../../stores/auth.ts'
import { handleError } from './handleAuthErrorUtil'
import { isEmailValid, isPasswordValid, isSixDigitCodeValid } from '../../utils/validationUtils';

enum State {
  Credentials,
  Mfa,
}

const router = useRouter()
const authStore = useAuthStore()

const email = ref<string>('')
const password = ref<string>('')
const otp = ref<string>('')

const isLoading = ref<boolean>(false)
const isLoadingMfa = ref<boolean>(false)
const serverError = ref<string>('')
const validationErrors = ref<Record<string, string>>({})
const state = ref<State>(State.Credentials)

const countdown = ref<number>(0)
let countdownTimer: number | undefined = undefined

const inputClasses: string = 'block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600 sm:text-sm/6'

const headerText = computed<string>(() => {
  switch (state.value) {
    case State.Credentials:
      return 'Sign in to Your Account'
    case State.Mfa:
      return 'Verify Your Email Address'
  }
})

const isCredentialsFormValid = computed<boolean>(() => {
  return isEmailValid(email.value) && isPasswordValid(password.value)
})

const isMfaFormValid = computed<boolean>(() => {
  return isSixDigitCodeValid(otp.value)
})

const validateField = (fieldName: string): void => {
  serverError.value = ''

  switch (fieldName) {
    case 'email':
      if (!email.value.trim()) {
        validationErrors.value.email = 'Email address is required.'
      } else if (!isEmailValid(email.value)) {
        validationErrors.value.email = 'Please enter a valid email address.'
      } else {
        delete validationErrors.value.email
      }
      break
    case 'password':
      if (!password.value) {
        validationErrors.value.password = 'Password is required.'
      } else if (!isPasswordValid(password.value)) {
        validationErrors.value.password =
          'Password must be between 6 and 100 characters long, and contain at least one uppercase letter, one lowercase letter, one digit, and one special character (e.g., !@#$%^&*).'
      } else {
        delete validationErrors.value.password
      }
      break
    case 'otp':
      if (!otp.value.trim()) {
        validationErrors.value.otp = 'Verification code is required.'
      } else if (!isSixDigitCodeValid(otp.value)) {
        validationErrors.value.otp = 'Verification code must be a 6-digit number.'
      } else {
        delete validationErrors.value.otp
      }
      break
  }
}

const validateCredentialsForm = (): boolean => {
  validationErrors.value = {}
  let isValid = true

  validateField('email');
  validateField('password');

  if (Object.keys(validationErrors.value).length > 0) {
    isValid = false;
  }
  return isValid
}

const validateMfaForm = (): boolean => {
  validationErrors.value = {}
  let isValid = true
  validateField('otp');

  if (Object.keys(validationErrors.value).length > 0) {
    isValid = false;
  }
  return isValid
}

const clearMessages = (): void => {
  serverError.value = ''
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
  otp.value = ''
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
    const response = await authService.signInPassword(email.value, password.value)

    if (response.data.isMfaRequired) {
      state.value = State.Mfa
      startCountdown()
    } else {
      authStore.setSession(
        response.data.accessTokenExpiry,
        email.value
      )
      setTimeout(() => {
        router.push('/dashboard')
      }, 1500)
    }
  } catch (err: any) {
    handleError(err, serverError)
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
    const response = await authService.verifyMfaSignIn(email.value, otp.value)

    authStore.setSession(
      response.data.accessTokenExpiry,
      email.value
    )
    setTimeout(() => {
      router.push('/dashboard')
    }, 1500)
  } catch (err: any) {
    handleError(err, serverError)
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
    await authService.signInPassword(email.value, password.value)
    startCountdown()
  } catch (err: any) {
    handleError(err, serverError)
  } finally {
    isLoading.value = false
  }
}

onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>

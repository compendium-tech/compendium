<template>
  <div class="flex min-h-screen flex-col items-center justify-center px-6 py-12 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-sm">
      <RouterLink to="/">
        <img src="/icon.svg" alt="Your Company" class="mx-auto h-20 w-auto" />
      </RouterLink>
      <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
        {{ headerText }}
      </h2>
    </div>

    <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
      <form v-if="currentStep === 'form'" @submit.prevent="handleSubmit" class="space-y-6">
        <div v-if="mode === 'signup'">
          <label for="name" class="block text-sm/6 font-medium text-gray-900">Full Name</label>
          <div class="mt-2">
            <input id="name" type="text" v-model.trim="name" required autocomplete="name" placeholder="John Doe"
              :class="inputClasses" />
            <p v-if="validationErrors.name" class="mt-2 text-sm text-red-600">{{ validationErrors.name }}</p>
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm/6 font-medium text-gray-900">Email address</label>
          <div class="mt-2">
            <input id="email" type="email" v-model.trim="email" required autocomplete="email"
              placeholder="johndoe@gmail.com" :class="inputClasses" />
            <p v-if="validationErrors.email" class="mt-2 text-sm text-red-600">{{ validationErrors.email }}</p>
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
          <div class="mt-2">
            <input id="password" type="password" v-model="password" required
              :autocomplete="mode === 'signup' ? 'new-password' : 'current-password'" placeholder="A strong password"
              :class="inputClasses" />
            <p v-if="validationErrors.password" class="mt-2 text-sm text-red-600">{{ validationErrors.password }}</p>
          </div>
        </div>

        <div>
          <button type="submit" :disabled="isLoading || !isFormValid"
            class="flex w-full justify-center rounded-md bg-primary-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-primary-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:bg-primary-400 disabled:cursor-not-allowed">
            {{ mode === 'signup' ? 'Sign Up' : 'Sign In' }}
            <span v-if="isLoading"
              class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent"></span>
          </button>
        </div>
        <p v-if="serverError" class="text-red-600 text-center text-sm mt-4">{{ serverError }}</p>
      </form>

      <div v-else-if="currentStep === 'verifyOtp'">
        <p class="mt-4 text-sm/6 text-gray-600">A One-Time Password has been sent to <span
            class="font-medium text-primary-600">{{ email }}</span>.</p>
        <form @submit.prevent="verifyMfa" class="space-y-6 mt-6">
          <div>
            <div class="mt-2">
              <input id="otp" type="text" v-model.trim="otp" placeholder="Enter 6-digit verification code" required
                maxlength="6" :class="inputClasses" />
              <p v-if="validationErrors.otp" class="mt-2 text-sm text-red-600">{{ validationErrors.otp }}</p>
            </div>
          </div>
          <div>
            <button type="submit" :disabled="isLoadingMfa || !isOtpFormValid"
              class="flex w-full justify-center rounded-md bg-primary-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-primary-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:bg-primary-400 disabled:cursor-not-allowed">
              Verify Account
              <span v-if="isLoadingMfa"
                class="ml-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-t-2 border-white border-t-transparent"></span>
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
      </div>

      <p v-if="currentStep === 'form'" class="mt-10 text-center text-sm/6 text-gray-500">
        <template v-if="mode === 'signup'">
          Already a member?
          <RouterLink to="/auth/signin" class="font-semibold text-primary-600 hover:text-primary-500">Sign in to your
            account</RouterLink>
        </template>
        <template v-else>
          Not a member?
          <RouterLink to="/auth/signup" class="font-semibold text-primary-600 hover:text-primary-500">Create your
            account</RouterLink>
        </template>
      </p>
    </div>
  </div>
</template>

<script>
import { RouterLink } from 'vue-router';
import { authService } from '../../api';
import { useAuthStore } from '../../stores/auth';

export default {
  components: { RouterLink },
  props: {
    mode: {
      type: String,
      default: 'signup', // Can be 'signup' or 'signin'
      validator: (value) => ['signup', 'signin'].includes(value),
    },
  },
  data() {
    return {
      name: '',
      email: '',
      password: '',
      otp: '',
      isLoading: false,
      isLoadingMfa: false,
      serverError: '',
      validationErrors: {},
      currentStep: 'form', // Can be 'form' (for signup/signin) or 'verifyOtp'
      countdown: 0,
      countdownTimer: null,
    };
  },
  computed: {
    inputClasses() {
      return 'block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-primary-600 sm:text-sm/6';
    },
    headerText() {
      if (this.currentStep === 'verifyOtp') {
        return 'Verify Your Email (OTP)';
      }
      return this.mode === 'signup' ? 'Create Your Account' : 'Log in to Your Account';
    },
    isFormValid() {
      if (this.mode === 'signup') {
        return this.name.length > 0 && this.validateEmail(this.email) && this.password.length >= 8;
      } else { // signin mode
        return this.validateEmail(this.email) && this.password.length >= 8;
      }
    },
    isOtpFormValid() {
      return /^\d{6}$/.test(this.otp);
    }
  },
  watch: {
    name() {
      if (this.validationErrors.name) delete this.validationErrors.name;
    },
    email() {
      if (this.validationErrors.email) delete this.validationErrors.email;
    },
    password() {
      if (this.validationErrors.password) delete this.validationErrors.password;
    },
    otp() {
      if (this.validationErrors.otp) delete this.validationErrors.otp;
    }
  },
  methods: {
    validateEmail(email) {
      return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    },
    validateForm() {
      this.validationErrors = {};
      let isValid = true;

      if (this.mode === 'signup' && !this.name.trim()) {
        this.validationErrors.name = 'Full Name is required.';
        isValid = false;
      }
      if (!this.email.trim()) {
        this.validationErrors.email = 'Email address is required.';
        isValid = false;
      } else if (!this.validateEmail(this.email)) {
        this.validationErrors.email = 'Please enter a valid email address.';
        isValid = false;
      }
      if (!this.password) {
        this.validationErrors.password = 'Password is required.';
        isValid = false;
      } else if (this.password.length < 8) {
        this.validationErrors.password = 'Password must be at least 8 characters long.';
        isValid = false;
      }
      return isValid;
    },
    validateOtpForm() {
      this.validationErrors = {};
      let isValid = true;
      if (!this.otp.trim()) {
        this.validationErrors.otp = 'OTP is required.';
        isValid = false;
      } else if (!/^\d{6}$/.test(this.otp)) {
        this.validationErrors.otp = 'OTP must be a 6-digit number.';
        isValid = false;
      }
      return isValid;
    },

    handleApiError(err) {
      this.serverError = '';

      if (err.response) {
        const statusCode = err.response.status;
        const errorMessage = err.response.data?.message || 'An unexpected error occurred.';

        switch (statusCode) {
          case 401:
            this.serverError = 'Invalid credentials or OTP. Please check and try again.';
            break;
          case 409:
            this.serverError = 'This email address is already registered. Please try logging in.';
            break;
          case 429:
            this.serverError = 'Too many requests. Please wait a moment before trying again.';
            break;
          case 500:
            this.serverError = 'Something went wrong on our end. Please try again later.';
            break;
          case 400:
            this.serverError = errorMessage;
            break;
          default:
            this.serverError = errorMessage;
            break;
        }
      } else if (err.request) {
        this.serverError = 'Network error. Please check your internet connection.';
      } else {
        this.serverError = err.message || 'An unexpected error occurred.';
      }
      console.error('API Error:', err);
    },

    clearMessages() {
      this.serverError = '';
      this.validationErrors = {};
    },

    async handleSubmit() {
      this.clearMessages();
      if (!this.validateForm()) {
        return;
      }

      this.isLoading = true;
      try {
        if (this.mode === 'signup') {
          await authService.signUp(this.name, this.email, this.password);
          this.currentStep = 'verifyOtp';
          this.startCountdown();
        } else {
          const response = await authService.signInPassword(this.email, this.password);
          if (response.data.isMfaRequired) {
            this.currentStep = 'verifyOtp';
            this.startCountdown();
          } else {
            const authStore = useAuthStore();
            authStore.setSession(
              response.data.accessTokenExpiry,
              this.email
            );
            setTimeout(() => {
              this.$router.push('/dashboard');
            }, 1500);
          }
        }
      } catch (err) {
        this.handleApiError(err);
      } finally {
        this.isLoading = false;
      }
    },

    async verifyMfa() {
      this.clearMessages();
      if (!this.validateOtpForm()) {
        return;
      }

      this.isLoadingMfa = true;
      const authStore = useAuthStore();
      try {
        const response = await authService.verifyMfaSignIn(this.email, this.otp); // Changed to verifyMfaSignIn
        authStore.setSession(
          response.data.accessTokenExpiry,
          this.email
        );
        setTimeout(() => {
          this.$router.push('/dashboard');
        }, 1500);
      } catch (err) {
        this.handleApiError(err);
      } finally {
        this.isLoadingMfa = false;
      }
    },

    async resendOtp() {
      this.clearMessages();
      if (this.countdown > 0 || this.isLoading) {
        return;
      }

      this.isLoading = true;
      try {
        if (this.mode === 'signup') {
          await authService.signUp(this.name, this.email, this.password);
        } else {
          await authService.signInPassword(this.email, this.password);
        }
        this.startCountdown();
      } catch (err) {
        this.handleApiError(err);
      } finally {
        this.isLoading = false;
      }
    },

    startCountdown() {
      this.countdown = 60;
      if (this.countdownTimer) {
        clearInterval(this.countdownTimer);
      }
      this.countdownTimer = setInterval(() => {
        if (this.countdown > 0) {
          this.countdown--;
        } else {
          clearInterval(this.countdownTimer);
          this.countdownTimer = null;
        }
      }, 1000);
    },

    goBackToForm() {
      this.currentStep = 'form';
      this.otp = '';
      this.clearMessages();
      clearInterval(this.countdownTimer);
      this.countdown = 0;
    },
  },
  beforeUnmount() {
    if (this.countdownTimer) {
      clearInterval(this.countdownTimer);
    }
  }
};
</script>

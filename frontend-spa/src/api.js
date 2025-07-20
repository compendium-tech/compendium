import axios from "axios"
import Cookies from "js-cookie"
import { useAuthStore } from "./stores/auth.js"

const API_BASE_URL = "http://localhost:8080/api/v1"

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

apiClient.interceptors.request.use(
  (config) => {
    const csrfToken = Cookies.get("csrfToken")
    if (csrfToken) {
      config.headers["X-CSRF-Token"] = csrfToken
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

export const AppErrorKind = {
  InternalServerError: 0,
  RequestValidationError: 1,
  InvalidCredentialsError: 2,
  EmailTakenError: 3,
  UserNotFoundError: 4,
  TooManyRequestsError: 5,
  MfaNotRequestedError: 6,
  InvalidMfaOtpError: 7,
  InvalidSessionError: 8,
}

export const authService = {
  signUp: (name, email, password) => {
    return apiClient.post("/users", { name, email, password })
  },
  verifyMfaSignUp: (email, otp) => {
    return apiClient.post("/sessions?flow=mfa", { email, otp })
  },

  signInPassword: (email, password) => {
    return apiClient.post("/sessions?flow=password", { email, password })
  },
  verifyMfaSignIn: (email, otp) => {
    return apiClient.post("/sessions?flow=mfa", { email, otp })
  },

  initiatePasswordReset: (email) => {
    return apiClient.put("/password?flow=init", { email })
  },
  confirmPasswordReset: (email, otp, newPassword) => {
    return apiClient.put("/password?flow=finish", {
      email,
      otp,
      password: newPassword,
    })
  },

  refresh: () => {
    return apiClient.post("/sessions?flow=refresh")
  },
}

export const userService = {
  getAccountDetails: async () => {
    return await apiClient.get("/account");
  },
  updateName: async (name) => {
    return await apiClient.put("/account", { name });
  },
}

let failedQueue = []

const processQueue = (error, tokenRefreshed = false) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else {
      prom.resolve(tokenRefreshed)
    }
  })
  failedQueue = []
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (window.location.pathname.startsWith("/auth/")) {
      return Promise.reject(error)
    }

    const originalRequest = error.config
    const authStore = useAuthStore()

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      if (authStore.isRefreshingToken) {
        return new Promise(function (resolve, reject) {
          failedQueue.push({
            resolve: () => resolve(apiClient(originalRequest)),
            reject,
          })
        })
      }

      console.log("Need to refresh")

      const refreshed = await authStore.refresh()

      if (refreshed) {
        processQueue(null, true)
        return apiClient(originalRequest)
      } else {
        processQueue(error)
        authStore.clearSession()
        return Promise.reject(error)
      }
    }

    return Promise.reject(error)
  }
)

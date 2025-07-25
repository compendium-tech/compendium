import { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from "axios"
import { apiClient, handleAxiosError } from "./base"
import { useAuthStore } from "../stores/auth"
import Cookies from "js-cookie"

interface SessionResponse {
  accessTokenExpiresAt: string
  refreshTokenExpiresAt: string
}

interface SignInResponse {
  isMfaRequired: boolean
  accessTokenExpiresAt?: string
  refreshTokenExpiresAt?: string
}

interface AuthService {
  signUp: (name: string, email: string, password: string) => Promise<void>
  verifyMfaSignUp: (email: string, otp: string) => Promise<SessionResponse>
  signInPassword: (email: string, password: string) => Promise<SignInResponse>
  verifyMfaSignIn: (email: string, otp: string) => Promise<SessionResponse>
  initiatePasswordReset: (email: string) => Promise<void>
  confirmPasswordReset: (
    email: string,
    otp: string,
    newPassword: string
  ) => Promise<any>
  refresh: () => Promise<SessionResponse>
}

export const authService: AuthService = {
  signUp: async (name, email, password) => {
    try {
      const response = await apiClient.post("/users", { name, email, password })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  verifyMfaSignUp: async (email, otp) => {
    try {
      const response = await apiClient.post("/sessions?flow=mfa", {
        email,
        otp,
      })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  signInPassword: async (email, password) => {
    try {
      const response = await apiClient.post("/sessions?flow=password", {
        email,
        password,
      })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  verifyMfaSignIn: async (email, otp) => {
    try {
      const response = await apiClient.post("/sessions?flow=mfa", {
        email,
        otp,
      })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  initiatePasswordReset: async (email) => {
    try {
      const response = await apiClient.put("/password?flow=init", { email })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  confirmPasswordReset: async (email, otp, newPassword) => {
    try {
      const response = await apiClient.put("/password?flow=finish", {
        email,
        otp,
        password: newPassword,
      })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  refresh: async () => {
    try {
      const response = await apiClient.post("/sessions?flow=refresh")

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
}

apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const csrfToken = Cookies.get("csrfToken")
    if (csrfToken) {
      config.headers.set("X-Csrf-Token", csrfToken)
    }

    const authStore = useAuthStore()
    const now = new Date()
    const expiryThreshold = now.getTime() + 5 * 1000

    const isAuthPath = window.location.pathname.startsWith("/auth/")
    const accessTokenExpiresAt = new Date(authStore.accessTokenExpiresAt)

    if (
      accessTokenExpiresAt &&
      accessTokenExpiresAt.getTime() < expiryThreshold &&
      !authStore.isRefreshingToken &&
      !isAuthPath
    ) {
      authStore.setIsRefreshingToken(true)
      try {
        console.log(
          "Access token nearing expiry, attempting to refresh proactively..."
        )
        const refreshResponse = await authService.refresh()
        authStore.setSession(refreshResponse.accessTokenExpiresAt)

        const csrfToken = Cookies.get("csrfToken")
        if (csrfToken) {
          config.headers.set("X-Csrf-Token", csrfToken)
        }
      } catch (refreshError) {
        authStore.clearSession()
        location.href = "/auth/signin"

        return Promise.reject(
          new Error("Session expired. Please sign in again.")
        )
      } finally {
        authStore.setIsRefreshingToken(false)
      }
    }

    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const authStore = useAuthStore()

    if (
      error.response?.status === 401 &&
      !window.location.pathname.startsWith("/auth/")
    ) {
      authStore.clearSession()

      location.href = "/auth/signin"
      return Promise.reject(error)
    }

    return Promise.reject(error)
  }
)

export interface Session {
  id: string
  isCurrent: boolean
  name: string
  os: string
  device: string
  location: string
  userAgent: string
  ipAddress: string
  createdAt: string
}

export interface SessionService {
  getSessions: () => Promise<Session[]>
  deleteSession: (id: string) => Promise<void>
  logout: () => Promise<void>
}

export const sessionService: SessionService = {
  getSessions: async () => {
    try {
      const response = await apiClient.get("/sessions")

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },

  deleteSession: async (id: string) => {
    try {
      await apiClient.delete(`/sessions/${id}`)
    } catch (error) {
      return handleAxiosError(error)
    }
  },

  logout: async () => {
    try {
      await apiClient.delete("/session")
    } catch (error) {
      return handleAxiosError(error)
    }
  },
}

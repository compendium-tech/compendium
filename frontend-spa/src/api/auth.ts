import { InternalAxiosRequestConfig } from "axios"
import { apiClient, ApiError, ApiErrorType } from "./base"
import { useAuthStore } from "../stores/auth"
import Cookies from "js-cookie"

export interface SessionResponse {
  accessTokenExpiresAt: string
  refreshTokenExpiresAt: string
}

export interface SignInResponse {
  isMfaRequired: boolean
  accessTokenExpiresAt?: string
  refreshTokenExpiresAt?: string
}

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
  ) => Promise<void>
  refresh: () => Promise<SessionResponse>
}

interface SessionService {
  getSessions: () => Promise<Session[]>
  deleteSession: (id: string) => Promise<void>
  logout: () => Promise<void>
}

export const authService: AuthService = {
  signUp: async (name, email, password) => {
    return (await apiClient.post("/users", { name, email, password })).data
  },
  verifyMfaSignUp: async (email, otp) => {
    return (
      await apiClient.post("/sessions?flow=mfa", {
        email,
        otp,
      })
    ).data
  },
  signInPassword: async (email, password) => {
    return (
      await apiClient.post("/sessions?flow=password", {
        email,
        password,
      })
    ).data
  },
  verifyMfaSignIn: async (email, otp) => {
    return (
      await apiClient.post("/sessions?flow=mfa", {
        email,
        otp,
      })
    ).data
  },
  initiatePasswordReset: async (email) => {
    await apiClient.put("/password?flow=init", { email })
  },
  confirmPasswordReset: async (email, otp, newPassword) => {
    await apiClient.put("/password?flow=finish", {
      email,
      otp,
      password: newPassword,
    })
  },
  refresh: async () => {
    return (await apiClient.post("/sessions?flow=refresh")).data
  },
}

export const sessionService: SessionService = {
  getSessions: async () => {
    return (await apiClient.get("/sessions")).data
  },

  deleteSession: async (id: string) => {
    await apiClient.delete(`/sessions/${id}`)
  },

  logout: async () => {
    await apiClient.delete("/session")
  },
}

/**
 * Interceptor to handle CSRF token and refresh the access token when it's about to expire.
 */
const authFlowInterceptor = async (config: InternalAxiosRequestConfig) => {
  // Add CSRF token to the request headers if it exists in cookies
  const csrfToken = Cookies.get("csrfToken")
  if (csrfToken) {
    config.headers.set("X-Csrf-Token", csrfToken)
  }

  const authStore = useAuthStore()
  const now = new Date()
  // Define a threshold (5 seconds) before the actual expiry to refresh the token
  const expiryThreshold = now.getTime() + 5 * 1000

  const isAuthPath = window.location.pathname.startsWith("/auth/")
  const accessTokenExpiresAt = new Date(authStore.accessTokenExpiresAt)

  // Check if the access token is about to expire, if we're not already refreshing the token, and if we're not on an authentication page
  if (
    accessTokenExpiresAt &&
    accessTokenExpiresAt.getTime() < expiryThreshold &&
    !authStore.isRefreshingToken &&
    !isAuthPath
  ) {
    // Set the refreshing token flag to prevent concurrent refresh attempts
    authStore.setIsRefreshingToken(true)
    try {
      console.log(
        "Access token nearing expiry, attempting to refresh proactively..."
      )
      // Call the refresh endpoint to get a new access token
      const refreshResponse = await authService.refresh()
      // Update the session with the new access token expiry time
      authStore.setSession(refreshResponse.accessTokenExpiresAt)

      // After refreshing the token, re-apply the CSRF token to the headers
      const csrfToken = Cookies.get("csrfToken")
      if (csrfToken) {
        config.headers.set("X-Csrf-Token", csrfToken)
      }
    } catch (refreshError) {
      // If refresh fails, clear the session and redirect to the sign-in page
      authStore.clearSession()
      location.href = "/auth/signin"

      return Promise.reject(new Error("Session expired. Please sign in again."))
    } finally {
      // Reset the refreshing token flag
      authStore.setIsRefreshingToken(false)
    }
  }

  return config
}

/**
 * Interceptor to redirect to the sign-in page if the API returns an invalid session error.
 */
const redirectToSignInInterceptor = (error: ApiError) => {
  const authStore = useAuthStore()

  // Check if the error is an invalid session error and if we're not already on an authentication page
  if (
    error.type === ApiErrorType.InvalidSessionError &&
    !window.location.pathname.startsWith("/auth/")
  ) {
    // Clear the session and redirect to the sign-in page
    authStore.clearSession()

    location.href = "/auth/signin"
    return Promise.reject(error)
  }

  return Promise.reject(error)
}

apiClient.interceptors.request.use(authFlowInterceptor, (error) =>
  Promise.reject(error)
)
apiClient.interceptors.response.use(
  (response) => response,
  redirectToSignInInterceptor
)

import axios, {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios"
import Cookies from "js-cookie"
import { useAuthStore } from "./stores/auth"

const API_BASE_URL = "http://localhost:1000/api/v1"

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

export enum ApiErrorKind {
  InternalServerError = 0,
  RequestValidationError = 1,
  InvalidCredentialsError = 2,
  EmailTakenError = 3,
  UserNotFoundError = 4,
  TooManyRequestsError = 5,
  MfaNotRequestedError = 6,
  InvalidMfaOtpError = 7,
  InvalidSessionError = 8,
}

export class ApiError extends Error {
  kind: ApiErrorKind

  constructor(kind: ApiErrorKind, message: string) {
    super(message)
    this.name = "ApiError"
    this.kind = kind
    Object.setPrototypeOf(this, ApiError.prototype)
  }
}

function parseAndRejectAxiosError(error: any): Promise<never> {
  if (error.response) {
    const { data } = error.response

    if (
      data &&
      typeof data === "object" &&
      "errorKind" in data &&
      "errorMessage" in data
    ) {
      return Promise.reject(
        new ApiError(
          (data as any).errorKind as ApiErrorKind,
          (data as any).errorMessage as string
        )
      )
    } else {
      return Promise.reject(
        new ApiError(
          ApiErrorKind.InternalServerError,
          "An unexpected error occurred."
        )
      )
    }
  } else if (error.request) {
    return Promise.reject(
      new ApiError(
        ApiErrorKind.InternalServerError,
        "No response received from server. Please check your network connection."
      )
    )
  } else {
    return Promise.reject(
      new ApiError(
        ApiErrorKind.InternalServerError,
        "An unexpected error occurred."
      )
    )
  }
}

interface SessionResponse {
  accessTokenExpiresAt: Date
  refreshTokenExpiresAt: Date
}

interface SignInResponse {
  isMfaRequired: boolean
  accessTokenExpiresAt?: Date
  refreshTokenExpiresAt?: Date
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
      return parseAndRejectAxiosError(error)
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
      return parseAndRejectAxiosError(error)
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
      return parseAndRejectAxiosError(error)
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
      return parseAndRejectAxiosError(error)
    }
  },
  initiatePasswordReset: async (email) => {
    try {
      const response = await apiClient.put("/password?flow=init", { email })
      return response.data
    } catch (error) {
      return parseAndRejectAxiosError(error)
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
      return parseAndRejectAxiosError(error)
    }
  },
  refresh: async () => {
    try {
      const response = await apiClient.post("/sessions?flow=refresh")
      return response.data
    } catch (error) {
      return parseAndRejectAxiosError(error)
    }
  },
}

interface UserService {
  getAccountDetails: () => Promise<any>
  updateName: (name: string) => Promise<any>
}

export const userService: UserService = {
  getAccountDetails: async () => {
    try {
      const response = await apiClient.get("/account")
      return response.data
    } catch (error) {
      return parseAndRejectAxiosError(error)
    }
  },
  updateName: async (name) => {
    try {
      const response = await apiClient.put("/account", { name })
      return response.data
    } catch (error) {
      return parseAndRejectAxiosError(error)
    }
  },
}

let isRefreshing = false

apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const csrfToken: string | undefined = Cookies.get("csrfToken")
    if (csrfToken) {
      config.headers.set("X-Csrf-Token", csrfToken)
    }

    const authStore = useAuthStore()
    const now = new Date()
    const expiryThreshold = new Date(now.getTime() + 5 * 1000)

    const isAuthPath = window.location.pathname.startsWith("/auth/")

    if (
      authStore.accessTokenExpiresAt &&
      authStore.accessTokenExpiresAt < expiryThreshold &&
      !isRefreshing &&
      !isAuthPath
    ) {
      isRefreshing = true
      try {
        console.log(
          "Access token nearing expiry, attempting to refresh proactively..."
        )
        const refreshResponse = await authService.refresh()
        authStore.setSession(refreshResponse.accessTokenExpiresAt)
        console.log("Access token refreshed proactively.")
      } catch (refreshError) {
        console.error("Proactive token refresh failed:", refreshError)
        authStore.clearSession()
        location.href = "/auth/signin"
        return Promise.reject(
          new Error("Session expired. Please sign in again.")
        )
      } finally {
        isRefreshing = false
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
      console.log(
        "Received 401 on non-auth path, clearing session and redirecting."
      )
      authStore.clearSession()
      location.href = "/auth/signin"
      return Promise.reject(error)
    }

    return Promise.reject(error)
  }
)

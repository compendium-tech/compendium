import axios, {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios"
import Cookies from "js-cookie"
import { useAuthStore } from "./stores/auth"

const API_BASE_URL = "http://localhost:8080/api/v1"

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const csrfToken: string | undefined = Cookies.get("csrfToken")
    if (csrfToken) {
      config.headers.set("X-Csrf-Token", csrfToken)
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

/**
 * Parses an AxiosError and throws a custom ApiError.
 * This centralizes error handling for API responses.
 * @param error The AxiosError to parse.
 * @throws ApiError A custom error containing kind and message.
 */
function parseAndThrowAxiosError(error: any): never {
  if (error.response) {
    const { data } = error.response

    if (
      data &&
      typeof data === "object" &&
      "errorKind" in data &&
      "errorMessage" in data
    ) {
      throw new ApiError(
        (data as any).errorKind as ApiErrorKind,
        (data as any).errorMessage as string
      )
    } else {
      throw new ApiError(
        ApiErrorKind.InternalServerError,
        "An unexpected error occurred."
      )
    }
  } else if (error.request) {
    throw new ApiError(
      ApiErrorKind.InternalServerError,
      "No response received from server. Please check your network connection."
    )
  } else {
    throw new ApiError(
      ApiErrorKind.InternalServerError,
      "An unexpected error occurred."
    )
  }
}

interface SessionResponse {
  accessTokenExpiry: Date
  refreshTokenExpiry: Date
}

interface SignInResponse {
  isMfaRequired: boolean
  accessTokenExpiry?: Date
  refreshTokenExpiry?: Date
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

// Implementation of the authentication service
export const authService: AuthService = {
  signUp: async (name, email, password) => {
    try {
      const response = await apiClient.post("/users", { name, email, password })
      return response.data
    } catch (error) {
      parseAndThrowAxiosError(error)
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
      parseAndThrowAxiosError(error)
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
      parseAndThrowAxiosError(error)
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
      parseAndThrowAxiosError(error)
    }
  },
  initiatePasswordReset: async (email) => {
    try {
      const response = await apiClient.put("/password?flow=init", { email })
      return response.data
    } catch (error) {
      parseAndThrowAxiosError(error)
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
      parseAndThrowAxiosError(error)
    }
  },
  refresh: async () => {
    try {
      const response = await apiClient.post("/sessions?flow=refresh")
      return response.data
    } catch (error) {
      parseAndThrowAxiosError(error)
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
      parseAndThrowAxiosError(error)
    }
  },
  updateName: async (name) => {
    try {
      const response = await apiClient.put("/account", { name })
      return response.data
    } catch (error) {
      parseAndThrowAxiosError(error)
    }
  },
}

interface RetryAxiosRequestConfig extends AxiosRequestConfig {
  _retry?: boolean
}

interface FailedRequestPromise {
  resolve: (value: AxiosResponse | PromiseLike<AxiosResponse>) => void
  reject: (reason?: any) => void
  originalRequest: RetryAxiosRequestConfig
}

let failedQueue: FailedRequestPromise[] = []

/**
 * Processes the queue of failed requests after a token refresh.
 * @param error An ApiError if the refresh failed, or null if successful.
 */
const processQueue = (error: ApiError | null): void => {
  failedQueue.forEach(async (prom) => {
    if (error) {
      prom.reject(error)
    } else {
      try {
        const response = await apiClient(prom.originalRequest)
        prom.resolve(response)
      } catch (err) {
        prom.reject(err)
      }
    }
  })

  failedQueue = []
}

apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    // If the request is to an auth path, don't try to refresh
    if (window.location.pathname.startsWith("/auth/")) {
      return Promise.reject(error)
    }

    const originalRequest = error.config as RetryAxiosRequestConfig
    const authStore = useAuthStore()

    // Check if it's a 401 Unauthorized error and hasn't been retried yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      // If a token refresh is already in progress, queue the current request
      if (authStore.isRefreshingToken) {
        return new Promise<AxiosResponse>((resolve, reject) => {
          failedQueue.push({
            resolve,
            reject,
            originalRequest,
          })
        })
      }

      authStore.setIsRefreshingToken(true)

      try {
        const refreshResponse = await authService.refresh()

        authStore.setSession(refreshResponse.accessTokenExpiry, authStore.email)

        processQueue(null)

        return apiClient(originalRequest)
      } catch (refreshError: any) {
        processQueue(refreshError)
        authStore.clearSession()
        location.href = "/auth/signin"
        return Promise.reject(refreshError)
      } finally {
        authStore.setIsRefreshingToken(false)
      }
    }
  }
)

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

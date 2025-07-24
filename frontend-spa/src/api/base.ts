import axios, { AxiosInstance } from "axios"

const API_BASE_URL = "http://localhost/v1"

export const apiClient: AxiosInstance = axios.create({
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

export function handleAxiosError(error: any): Promise<never> {
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

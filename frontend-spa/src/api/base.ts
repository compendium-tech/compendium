import axios, { AxiosInstance } from "axios"

const API_BASE_URL = "http://localhost/v1"

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

export enum ApiErrorType {
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
  type: ApiErrorType

  constructor(type: ApiErrorType) {
    super("API error")
    this.name = "ApiError"
    this.type = type
    Object.setPrototypeOf(this, ApiError.prototype)
  }
}

export function handleAxiosError(error: any): Promise<never> {
  if (error.response) {
    const { data } = error.response

    if (
      data &&
      typeof data === "object" &&
      "errorType" in data
    ) {
      return Promise.reject(
        new ApiError(
          (data as any).errorType as ApiErrorType
        )
      )
    } else {
      return Promise.reject(
        new ApiError(
          ApiErrorType.InternalServerError,
        )
      )
    }
  } else if (error.request) {
    return Promise.reject(
      new ApiError(
        ApiErrorType.InternalServerError,
      )
    )
  } else {
    return Promise.reject(
      new ApiError(
        ApiErrorType.InternalServerError,
      )
    )
  }
}

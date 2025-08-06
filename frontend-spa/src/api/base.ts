import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios"

const API_BASE_URL = "http://localhost/v1"

export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
})

/**
 * Response interceptor for handling API errors.
 *
 * This interceptor checks the response for error conditions and rejects the promise
 * with an ApiError if a known error type is returned by the API.  If the error
 * is not a known type, or if the error is due to a network issue, it rejects
 * with a generic InternalServerError.
 */
const errorInterceptor = (error: AxiosError): Promise<void> => {
  if (error.response) {
    const { data } = error.response

    if (data && typeof data === "object" && "errorType" in data) {
      return Promise.reject(
        new ApiError((data as any).errorType as ApiErrorType)
      )
    } else {
      return Promise.reject(new ApiError(ApiErrorType.InternalServerError))
    }
  } else if (error.request) {
    return Promise.reject(new ApiError(ApiErrorType.InternalServerError))
  } else {
    return Promise.reject(new ApiError(ApiErrorType.InternalServerError))
  }
}

apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  errorInterceptor
)

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
    super(errorMessage(type))

    this.name = "ApiError"
    this.type = type
    Object.setPrototypeOf(this, ApiError.prototype)
  }
}

const errorMessage = (errorType: ApiErrorType): string => {
  switch (errorType) {
    case ApiErrorType.RequestValidationError:
      return "Invalid request data. Please check your input."
    case ApiErrorType.InvalidCredentialsError:
      return "Invalid email or password. Please try again."
    case ApiErrorType.EmailTakenError:
      return "This email address is already registered. Please try logging in."
    case ApiErrorType.UserNotFoundError:
      return "User not found. Please check your email address."
    case ApiErrorType.TooManyRequestsError:
      return "Too many requests. Please wait a moment before trying again."
    case ApiErrorType.MfaNotRequestedError:
      return "MFA was not requested for this session."
    case ApiErrorType.InvalidMfaOtpError:
      return "Invalid OTP. Please check the code and try again."
    case ApiErrorType.InvalidSessionError:
      return "Your session is invalid or expired. Please sign in again."
    default:
      return "An unexpected error occured."
  }
}

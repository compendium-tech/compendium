import { Ref } from "vue"
import { ApiError, ApiErrorType } from "../../api/base"

export const handleApiError = (
  error: ApiError,
  targetReactiveObject: Ref<string>
) => {
  targetReactiveObject.value = getErrorMessage(error)
}

const getErrorMessage = (error: ApiError) => {
  switch (error.type) {
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

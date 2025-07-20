import { Ref } from "vue"
import { ApiError, ApiErrorKind } from "../../api"

export const handleApiError = (
  error: ApiError,
  targetReactiveObject: Ref<string>
) => {
  targetReactiveObject.value = getErrorMessage(error)
}

const getErrorMessage = (error: ApiError) => {
  switch (error.kind) {
    case ApiErrorKind.RequestValidationError:
      return "Invalid request data. Please check your input."
    case ApiErrorKind.InvalidCredentialsError:
      return "Invalid email or password. Please try again."
    case ApiErrorKind.EmailTakenError:
      return "This email address is already registered. Please try logging in."
    case ApiErrorKind.UserNotFoundError:
      return "User not found. Please check your email address."
    case ApiErrorKind.TooManyRequestsError:
      return "Too many requests. Please wait a moment before trying again."
    case ApiErrorKind.MfaNotRequestedError:
      return "MFA was not requested for this session."
    case ApiErrorKind.InvalidMfaOtpError:
      return "Invalid OTP. Please check the code and try again."
    case ApiErrorKind.InvalidSessionError:
      return "Your session is invalid or expired. Please sign in again."
    default:
      return "An unexpected error occured."
  }
}

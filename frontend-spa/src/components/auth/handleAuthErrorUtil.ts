import { Ref } from "vue"
import { AppErrorKind } from "../../api.js"

export function handleError(err: any, targetReactiveObject: Ref<string>) {
  targetReactiveObject.value = ""

  if (err.response) {
    const errorKind = err.response.data?.errorKind
    const defaultErrorMessage = "An unexpected error occurred."

    switch (errorKind) {
      case AppErrorKind.RequestValidationError:
        targetReactiveObject.value =
          "Invalid request data. Please check your input."
        break
      case AppErrorKind.InvalidCredentialsError:
        targetReactiveObject.value =
          "Invalid email or password. Please try again."
        break
      case AppErrorKind.EmailTakenError:
        targetReactiveObject.value =
          "This email address is already registered. Please try logging in."
        break
      case AppErrorKind.UserNotFoundError:
        targetReactiveObject.value =
          "User not found. Please check your email address."
        break
      case AppErrorKind.TooManyRequestsError:
        targetReactiveObject.value =
          "Too many requests. Please wait a moment before trying again."
        break
      case AppErrorKind.MfaNotRequestedError:
        targetReactiveObject.value = "MFA was not requested for this session."
        break
      case AppErrorKind.InvalidMfaOtpError:
        targetReactiveObject.value =
          "Invalid OTP. Please check the code and try again."
        break
      case AppErrorKind.InvalidSessionError:
        targetReactiveObject.value =
          "Your session is invalid or expired. Please sign in again."
        break
      default:
        targetReactiveObject.value = defaultErrorMessage
        break
    }
  } else if (err.request) {
    targetReactiveObject.value =
      "Network error. Please check your internet connection."
  } else {
    targetReactiveObject.value = "Unexpected error happened."
  }
}

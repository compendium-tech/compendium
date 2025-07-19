export const isEmailValid = (emailString: string): boolean => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailString)
}

export const isPasswordValid = (passwordString: string): boolean => {
  if (passwordString.length < 6 || passwordString.length > 100) {
    return false
  }

  if (!/[A-Z]/.test(passwordString)) {
    return false
  }
  if (!/[a-z]/.test(passwordString)) {
    return false
  }
  if (!/[0-9]/.test(passwordString)) {
    return false
  }
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]/.test(passwordString)) {
    return false
  }

  return true
}

export const isSixDigitCodeValid = (otpString: string): boolean => {
  return /^\d{6}$/.test(otpString)
}

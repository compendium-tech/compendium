package email

type EmailSender interface {
	SendSignUpMfaEmail(to string, otp string) error
	SendSignInMfaEmail(to string, otp string) error
}

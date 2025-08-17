package repository

import (
	"context"
)

type MfaRepository interface {
	SetMfaOtpByEmail(ctx context.Context, email string, otp string)
	GetMfaOtpByEmail(ctx context.Context, email string) *string
	RemoveMfaOtpByEmail(ctx context.Context, email string)
}

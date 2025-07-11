package repository

import (
	"context"
)

type MfaRepository interface {
	SetMfaOtpByEmail(ctx context.Context, email string, otp string) error
	GetMfaOtpByEmail(ctx context.Context, email string) (*string, error)
	RemoveMfaOtpByEmail(ctx context.Context, email string) error
}

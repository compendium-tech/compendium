package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	mfaOtpKeyPrefix = "mfa_otp:"
	mfaOtpTtl       = 120 * time.Second
)

type redisMfaRepository struct {
	client *redis.Client
}

func NewRedisMfaRepository(client *redis.Client) MfaRepository {
	return &redisMfaRepository{client: client}
}

func (r *redisMfaRepository) SetMfaOtpByEmail(ctx context.Context, email string, otp string) {
	err := r.client.Set(ctx, r.createMfaOtpKey(email), otp, mfaOtpTtl).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisMfaRepository) GetMfaOtpByEmail(ctx context.Context, email string) *string {
	code, err := r.client.Get(ctx, r.createMfaOtpKey(email)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		panic(err)
	}

	return &code
}

func (r *redisMfaRepository) RemoveMfaOtpByEmail(ctx context.Context, email string) {
	err := r.client.Del(ctx, r.createMfaOtpKey(email)).Err()
	if err != nil {
		panic(err)
	}
}

func (r *redisMfaRepository) createMfaOtpKey(email string) string {
	return fmt.Sprintf("%s%s", mfaOtpKeyPrefix, email)
}

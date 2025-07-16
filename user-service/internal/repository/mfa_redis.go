package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ztrue/tracerr"
)

const (
	mfaOtpKeyPrefix = "mfa_otp:"
	mfaOtpTtl       = 120 * time.Second
)

type RedisMfaRepository struct {
	client *redis.Client
}

func NewRedisMfaRepository(client *redis.Client) *RedisMfaRepository {
	return &RedisMfaRepository{client: client}
}

func (r *RedisMfaRepository) SetMfaOtpByEmail(ctx context.Context, email string, otp string) error {
	err := r.client.Set(ctx, r.createMfaOtpKey(email), otp, mfaOtpTtl).Err()

	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *RedisMfaRepository) GetMfaOtpByEmail(ctx context.Context, email string) (*string, error) {
	code, err := r.client.Get(ctx, r.createMfaOtpKey(email)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, tracerr.Wrap(err)
	}

	return &code, nil
}

func (r *RedisMfaRepository) RemoveMfaOtpByEmail(ctx context.Context, email string) error {
	err := r.client.Del(ctx, r.createMfaOtpKey(email)).Err()

	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *RedisMfaRepository) createMfaOtpKey(email string) string {
	return fmt.Sprintf("%s%s", mfaOtpKeyPrefix, email)
}

package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type _userIdKey struct{}

var userIdKey = _userIdKey{}

func GetUserId(ctx context.Context) (uuid.UUID, error) {
	if userId, ok := ctx.Value(userIdKey).(uuid.UUID); ok {
		return userId, nil
	} else {
		return uuid.Nil, tracerr.New("Authentication middleware didn't set user id to UUID value, perhaps it wasn't enabled?")
	}
}

func SetUserId(ctx *context.Context, userId uuid.UUID) {
	*ctx = context.WithValue(*ctx, userIdKey, userId)
}

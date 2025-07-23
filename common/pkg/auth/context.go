package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type _userIDKey struct{}

var userIDKey = _userIDKey{}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if userID, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return userID, nil
	} else {
		return uuid.Nil, tracerr.New("Authentication middleware didn't set user id to UUID value, perhaps it wasn't enabled?")
	}
}

func SetUserID(ctx *context.Context, userID uuid.UUID) {
	*ctx = context.WithValue(*ctx, userIDKey, userID)
}

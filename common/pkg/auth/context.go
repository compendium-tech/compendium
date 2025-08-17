package auth

import (
	"context"

	"github.com/google/uuid"
)

type _userIDKey struct{}

var userIDKey = _userIDKey{}

func GetUserID(ctx context.Context) uuid.UUID {
	id := GetUserIDOrNil(ctx)
	if id == uuid.Nil {
		panic("Authentication middleware didn't set user id to UUID value, perhaps it wasn't enabled?")
	}

	return id
}

func GetUserIDOrNil(ctx context.Context) uuid.UUID {
	if userID, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return userID
	} else {
		return uuid.Nil
	}
}

func SetUserID(ctx *context.Context, userID uuid.UUID) {
	*ctx = context.WithValue(*ctx, userIDKey, userID)
}

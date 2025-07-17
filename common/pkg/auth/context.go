package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type userIdKey struct{}

var UserIdKey = userIdKey{}

func GetUserId(ctx context.Context) (uuid.UUID, error) {
	if userId, ok := ctx.Value(UserIdKey).(uuid.UUID); ok {
		return userId, nil
	} else {
		return uuid.Nil, tracerr.New("Authentication middleware didn't set user id to UUID value, perhaps it wasn't enabled?")
	}
}

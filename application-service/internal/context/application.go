package localcontext

import (
	"context"

	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/ztrue/tracerr"
)

type _applicationKey struct{}

var applicationKey = _applicationKey{}

func SetApplication(ctx *context.Context, application *model.Application) {
	*ctx = context.WithValue(*ctx, applicationKey, application)
}

func GetApplication(ctx context.Context) (*model.Application, error) {
	if application, ok := ctx.Value(applicationKey).(*model.Application); ok {
		return application, nil
	} else {
		return nil, tracerr.New("middleware didn't set current application value, perhaps it wasn't enabled?")
	}
}

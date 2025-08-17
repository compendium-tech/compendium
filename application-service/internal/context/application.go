package localcontext

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/application-service/internal/model"
)

type _applicationKey struct{}

var applicationKey = _applicationKey{}

func SetApplication(ctx *context.Context, application model.Application) {
	*ctx = context.WithValue(*ctx, applicationKey, application)
}

func GetApplication(ctx context.Context) model.Application {
	if application, ok := ctx.Value(applicationKey).(model.Application); ok {
		return application
	} else {
		panic(fmt.Errorf("middleware didn't set current application value, perhaps it wasn't enabled?"))
	}
}

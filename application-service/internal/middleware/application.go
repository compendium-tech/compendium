package middleware

import (
	"fmt"

	localctx "github.com/compendium-tech/compendium/application-service/internal/context"
	appErr "github.com/compendium-tech/compendium/application-service/internal/error"
	"github.com/compendium-tech/compendium/application-service/internal/service"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SetApplicationFromRequest struct {
	applicationService service.ApplicationService
}

func NewSetApplicationFromRequest(applicationService service.ApplicationService) *SetApplicationFromRequest {
	return &SetApplicationFromRequest{
		applicationService: applicationService,
	}
}

func (s *SetApplicationFromRequest) Handle(c *gin.Context) {
	appErr.Handle(s.handle)(c)
}

func (s *SetApplicationFromRequest) handle(c *gin.Context) error {
	applicationIdString := c.Request.PathValue("applicationId")

	if applicationIdString == "" {
		return appErr.NewWithReason(appErr.RequestValidationError, "missing application ID")
	}

	applicationId, err := uuid.Parse(applicationIdString)

	if err != nil {
		return appErr.NewWithReason(appErr.RequestValidationError, fmt.Sprintf("application ID is not a valid UUID: %v", err))
	}

	application, err := s.applicationService.GetCurrentApplicationModel(c.Request.Context(), applicationId)
	if err != nil {
		return err
	}

	ctx := c.Request.Context()
	log.SetLogger(&ctx, log.L(ctx).WithField("applicationId", applicationId))
	localctx.SetApplication(&ctx, application)
	c.Request = c.Request.WithContext(ctx)

	return nil
}

package middleware

import (
	"fmt"

	localcontext "github.com/compendium-tech/compendium/application-service/internal/context"
	myerror "github.com/compendium-tech/compendium/application-service/internal/error"
	"github.com/compendium-tech/compendium/application-service/internal/service"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"
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
	var eh httputils.ErrorHandler

	eh.Handle(s.handle)(c)
}

func (s *SetApplicationFromRequest) handle(c *gin.Context) {
	applicationIdString := c.Request.PathValue("applicationId")

	if applicationIdString == "" {
		myerror.NewWithReason(myerror.RequestValidationError, "missing application ID").Throw()
	}

	applicationId, err := uuid.Parse(applicationIdString)

	if err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("application ID is not a valid UUID: %v", err)).Throw()
	}

	application := s.applicationService.GetCurrentApplicationModel(c.Request.Context(), applicationId)

	ctx := c.Request.Context()
	log.SetLogger(&ctx, log.L(ctx).WithField("applicationId", applicationId))
	localcontext.SetApplication(&ctx, application)
	c.Request = c.Request.WithContext(ctx)
}

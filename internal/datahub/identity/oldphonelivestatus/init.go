package oldphonelivestatus

import (
	"front-office/configs/application"
	"front-office/internal/core/log/operation"
	"front-office/internal/core/member"
	"front-office/internal/core/role"
	"front-office/internal/middleware"
	"front-office/pkg/httpclient"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *application.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client, nil)
	memberRepository := member.NewRepository(cfg, client, nil)
	roleRepository := role.NewRepository(cfg, client)
	logOperationRepo := operation.NewRepository(cfg, client, nil)
	service := NewService(repository, memberRepository)
	memberService := member.NewService(memberRepository, roleRepository, logOperationRepo)
	controller := NewController(service, memberService)

	phoneLiveStatusGroup := apiGroup.Group("old-phone-live-status")
	phoneLiveStatusGroup.Get("/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	phoneLiveStatusGroup.Get("/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
	phoneLiveStatusGroup.Get("/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	phoneLiveStatusGroup.Get("/jobs/:id/details", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	phoneLiveStatusGroup.Get("/jobs/:id/details/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetails)
	phoneLiveStatusGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(phoneLiveStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.SingleSearch)
	phoneLiveStatusGroup.Post("/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkSearch)
}

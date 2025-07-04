package phonelivestatus

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client, nil)
	memberRepository := member.NewRepository(cfg, client)
	roleRepository := role.NewRepository(cfg, client)
	logOperationRepo := operation.NewRepository(cfg, client)
	service := NewService(repository, memberRepository)
	memberService := member.NewService(memberRepository, roleRepository, logOperationRepo)
	controller := NewController(service, memberService)

	phoneLiveStatusGroup := apiGroup.Group("phone-live-status")
	phoneLiveStatusGroup.Get("/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	phoneLiveStatusGroup.Get("/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
	phoneLiveStatusGroup.Get("/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	phoneLiveStatusGroup.Get("/jobs/:id/details", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	phoneLiveStatusGroup.Get("/jobs/:id/details/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetails)
	phoneLiveStatusGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(phoneLiveStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.SingleSearch)
	phoneLiveStatusGroup.Post("/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkSearch)
}

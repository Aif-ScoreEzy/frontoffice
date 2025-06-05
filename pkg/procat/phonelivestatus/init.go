package phonelivestatus

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config) {
	repository := NewRepository(cfg)
	service := NewService(repository)
	controller := NewController(service)

	phoneLiveStatusGroup := apiGroup.Group("phone-live-status")
	phoneLiveStatusGroup.Get("/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	phoneLiveStatusGroup.Get("/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
	phoneLiveStatusGroup.Get("/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	phoneLiveStatusGroup.Get("/jobs/:id/details", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	phoneLiveStatusGroup.Get("/jobs/:id/details/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetails)
	phoneLiveStatusGroup.Post("/single-search", middleware.Auth(), middleware.IsRequestValid(PhoneLiveStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.SingleSearch)
	phoneLiveStatusGroup.Post("/bulk-search", middleware.Auth(), middleware.IsRequestValid(PhoneLiveStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.BulkSearch)
}

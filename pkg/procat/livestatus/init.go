package livestatus

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(liveStatusAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(cfg, repository)
	controller := NewController(service)

	liveStatusAPI.Post("/live-status", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.Search)
	liveStatusAPI.Post("/live-status/bulk", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.UploadCSVFile(), controller.BulkSearch)
	liveStatusAPI.Get("/live-status", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	liveStatusAPI.Get("/live-status/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	liveStatusAPI.Get("/live-status/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
	liveStatusAPI.Get("/live-status/:id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	liveStatusAPI.Get("/live-status/:id/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetailsExport)

	// Cron Reprocess Unsuccessful Job Details
	// jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	// scd := gocron.NewScheduler(jakartaTime)
	// _, _ = scd.Every(5).Minute().Do(controller.ReprocessFailedJobDetails)
	// scd.StartAsync()
}

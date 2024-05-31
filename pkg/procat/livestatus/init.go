package livestatus

import (
	"front-office/app/config"
	"front-office/pkg/middleware"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(liveStatusAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(cfg, repository)
	controller := NewController(service)

	liveStatusAPI.Post("/live-status", middleware.Auth(), middleware.GetPayloadFromJWT(), middleware.UploadCSVFile(), controller.BulkSearch)
	liveStatusAPI.Get("/live-status", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.GetJobs)
	liveStatusAPI.Get("/live-status/jobs-summary", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.GetJobsSummary)
	liveStatusAPI.Get("/live-status/jobs-summary/export", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.ExportJobsSummary)
	liveStatusAPI.Get("/live-status/:id", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.GetJobDetails)
	liveStatusAPI.Get("/live-status/:id/export", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.GetJobDetailsExport)

	// Cron Reprocess Unsuccessful Job Details
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scd := gocron.NewScheduler(jakartaTime)
	_, _ = scd.Every(5).Minute().Do(controller.ReprocessFailedJobDetails)
	scd.StartAsync()
}

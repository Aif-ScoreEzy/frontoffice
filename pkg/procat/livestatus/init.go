package livestatus

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(apiGroup fiber.Router, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(cfg, repository)
	controller := NewController(service)

	liveStatusGroup := apiGroup.Group("live-status")
	liveStatusGroup.Post("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.Search)
	liveStatusGroup.Post("/bulk", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.UploadCSVFile(), controller.BulkSearch)
	liveStatusGroup.Get("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	liveStatusGroup.Get("/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	liveStatusGroup.Get("/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
	liveStatusGroup.Get("/:id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	liveStatusGroup.Get("/:id/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetailsExport)

	// Cron Reprocess Unsuccessful Job Details
	// jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	// scd := gocron.NewScheduler(jakartaTime)
	// _, _ = scd.Every(5).Minute().Do(controller.ReprocessFailedJobDetails)
	// scd.StartAsync()
}

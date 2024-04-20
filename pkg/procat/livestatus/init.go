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

	liveStatusAPI.Post("/live-status", middleware.UploadCSVFile(), controller.BulkSearch)
	liveStatusAPI.Get("/live-status", middleware.Auth(), controller.GetJobs)
	liveStatusAPI.Get("/live-status/:id", middleware.Auth(), controller.GetJobDetails)

	// Cron Reprocess Unsuccessful Job Details
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scd := gocron.NewScheduler(jakartaTime)
	_, _ = scd.Every(5).Minute().Do(controller.ReprocessUnsuccessfulJobDetails)
	scd.StartAsync()
}

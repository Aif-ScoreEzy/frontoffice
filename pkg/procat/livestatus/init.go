package livestatus

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(liveStatusAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(repository)
	controller := NewController(service, cfg)

	liveStatusAPI.Post("/live-status", middleware.UploadCSVFile(), controller.BulkSearch)
	liveStatusAPI.Get("/live-status", middleware.Auth(), controller.GetJobs)
}

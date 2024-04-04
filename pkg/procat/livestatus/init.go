package livestatus

import (
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(liveStatusAPI fiber.Router, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	controller := NewController(service)

	liveStatusAPI.Post("/live-status", middleware.UploadCSVFile(), controller.UploadCSV)
}

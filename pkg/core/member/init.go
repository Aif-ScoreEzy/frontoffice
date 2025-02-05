package member

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(userAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	service := NewService(repo)
	controller := NewController(service)

	userAPI.Get("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetList)
	userAPI.Put("/profile", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(UpdateProfileRequest{}), controller.UpdateProfile)
	userAPI.Put("/upload-profile-image", middleware.Auth(), middleware.IsRequestValid(UploadProfileImageRequest{}), middleware.GetJWTPayloadFromCookie(), middleware.FileUpload(), controller.UploadProfileImage)
	userAPI.Get("/by", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetBy)
	userAPI.Get("/:id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetById)
	userAPI.Delete("/:id", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.DeleteById)
}

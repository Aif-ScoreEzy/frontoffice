package user

import (
	"front-office/app/config"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(userAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	repoRole := role.NewRepository(db)
	service := NewService(repo, repoRole)
	serviceRole := role.NewService(repoRole)
	controller := NewController(service, serviceRole)

	userAPI.Get("/", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), controller.GetAllUsers)
	userAPI.Put("/profile", middleware.Auth(), middleware.IsRequestValid(UpdateProfileRequest{}), middleware.GetPayloadFromJWT(), controller.UpdateProfile)
	userAPI.Put("/upload-profile-image", middleware.Auth(), middleware.IsRequestValid(UploadProfileImageRequest{}), middleware.GetPayloadFromJWT(), middleware.FileUpload(), controller.UploadProfileImage)
	userAPI.Get("/:id", middleware.Auth(), middleware.GetPayloadFromJWT(), controller.GetUserById)
	userAPI.Put("/:id", middleware.AdminAuth(), middleware.IsRequestValid(UpdateUserRequest{}), middleware.GetPayloadFromJWT(), controller.UpdateUserById)
	userAPI.Delete("/:id", middleware.AdminAuth(), middleware.GetPayloadFromJWT(), controller.DeleteUserById)
}

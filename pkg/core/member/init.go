package member

import (
	"front-office/app/config"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(userAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db, cfg)
	roleRepo := role.NewRepository(db, cfg)
	repoLogOperation := operation.NewRepository(cfg)

	serviceRole := role.NewService(roleRepo)
	service := NewService(repo, serviceRole)
	serviceLogOperation := operation.NewService(repoLogOperation)

	controller := NewController(service, serviceRole, serviceLogOperation)

	userAPI.Get("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetList)
	userAPI.Put("/profile", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(UpdateProfileRequest{}), controller.UpdateProfile)
	userAPI.Put("/upload-profile-image", middleware.Auth(), middleware.IsRequestValid(UploadProfileImageRequest{}), middleware.GetJWTPayloadFromCookie(), middleware.FileUpload(), controller.UploadProfileImage)
	userAPI.Get("/by", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetBy)
	userAPI.Get("/:id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetById)
	userAPI.Put("/:id", middleware.AdminAuth(), middleware.IsRequestValid(UpdateUserRequest{}), middleware.GetJWTPayloadFromCookie(), controller.UpdateMemberById)
	userAPI.Delete("/:id", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.DeleteById)
}

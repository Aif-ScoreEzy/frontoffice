package member

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/role"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(userAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client)
	roleRepo := role.NewRepository(cfg)
	repoLogOperation := operation.NewRepository(cfg, client)

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

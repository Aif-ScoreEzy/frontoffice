package permission

import (
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(permissionAPI fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	permissionAPI.Post("/", middleware.Auth(), middleware.IsRequestValid(PermissionRequest{}), controller.CreatePermission)
	permissionAPI.Get("/:id", middleware.Auth(), controller.GetPermissionById)
	permissionAPI.Put("/:id", middleware.Auth(), middleware.IsRequestValid(PermissionRequest{}), controller.UpdatePermissionById)
	permissionAPI.Delete("/:id", middleware.Auth(), controller.DeletePermissionById)
}

package role

import (
	"front-office/pkg/core/permission"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(roleAPI fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	repoPermission := permission.NewRepository(db)
	service := NewService(repo)
	servicePermission := permission.NewService(repoPermission)
	controller := NewController(service, servicePermission)

	roleAPI.Post("/", middleware.Auth(), middleware.IsRequestValid(CreateRoleRequest{}), controller.CreateRole)
	roleAPI.Get("/", middleware.Auth(), controller.GetAllRoles)
	roleAPI.Get("/:id", middleware.Auth(), controller.GetRoleById)
	roleAPI.Put("/:id", middleware.Auth(), middleware.IsRequestValid(UpdateRoleRequest{}), controller.UpdateRole)
	roleAPI.Delete("/:id", middleware.Auth(), controller.DeleteRole)
}

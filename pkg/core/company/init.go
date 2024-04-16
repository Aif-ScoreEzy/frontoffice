package company

import (
	"front-office/pkg/core/industry"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(companyAPI fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	repoIndustry := industry.NewRepository(db)
	service := NewService(repo)
	serviceIndustry := industry.NewService(repoIndustry)
	controller := NewController(service, serviceIndustry)

	companyAPI.Put("/:id", middleware.IsRequestValid(UpdateCompanyRequest{}), controller.UpdateCompanyByID)
}

package core

import (
	"front-office/app/config"
	"front-office/pkg/core/auth"
	"front-office/pkg/core/company"
	"front-office/pkg/core/grading"
	"front-office/pkg/core/log"
	"front-office/pkg/core/permission"
	"front-office/pkg/core/role"
	"front-office/pkg/core/user"
	"front-office/pkg/procat/livestatus"
	"front-office/pkg/scoreezy/genretail"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(routeAPI fiber.Router, cfg *config.Config, db *gorm.DB) {
	userAPI := routeAPI.Group("users")
	auth.SetupInit(userAPI, db, cfg)
	user.SetupInit(userAPI, db, cfg)

	roleAPI := routeAPI.Group("roles")
	role.SetupInit(roleAPI, db)

	permissionAPI := routeAPI.Group("permissions")
	permission.SetupInit(permissionAPI, db)

	companyAPI := routeAPI.Group("companies")
	company.SetupInit(companyAPI, db)

	gradingAPI := routeAPI.Group("gradings")
	grading.SetupInit(gradingAPI, db)

	genRetailAPI := routeAPI.Group("scores")
	genretail.SetupInit(genRetailAPI, db, cfg)

	logAPI := routeAPI.Group("logs")
	log.SetupInit(logAPI, cfg)

	productAPI := routeAPI.Group("products")
	livestatus.SetupInit(productAPI, db)
}

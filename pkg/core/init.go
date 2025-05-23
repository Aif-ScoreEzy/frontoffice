package core

import (
	"front-office/app/config"
	"front-office/pkg/core/auth"
	"front-office/pkg/core/company"
	"front-office/pkg/core/grading"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/member"
	"front-office/pkg/core/permission"
	"front-office/pkg/core/role"
	"front-office/pkg/procat/loanrecordchecker"
	"front-office/pkg/procat/multipleloan"
	"front-office/pkg/procat/phonelivestatus"
	"front-office/pkg/scoreezy/genretail"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(routeGroup fiber.Router, cfg *config.Config, db *gorm.DB) {
	userGroup := routeGroup.Group("users")
	auth.SetupInit(userGroup, db, cfg)
	member.SetupInit(userGroup, db, cfg)

	roleGroup := routeGroup.Group("roles")
	role.SetupInit(roleGroup, cfg, db)

	permissionGroup := routeGroup.Group("permissions")
	permission.SetupInit(permissionGroup, db)

	companyGroup := routeGroup.Group("companies")
	company.SetupInit(companyGroup, db)

	gradingGroup := routeGroup.Group("gradings")
	grading.SetupInit(gradingGroup, db, cfg)

	genRetailGroup := routeGroup.Group("scores")
	genretail.SetupInit(genRetailGroup, db, cfg)

	logGroup := routeGroup.Group("logs")
	transaction.SetupInit(logGroup, db, cfg)
	operation.SetupInit(logGroup, cfg)

	productGroup := routeGroup.Group("products")
	phonelivestatus.SetupInit(productGroup, cfg)
	loanrecordchecker.SetupInit(productGroup, cfg)
	multipleloan.SetupInit(productGroup, cfg)
}

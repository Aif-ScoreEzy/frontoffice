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
	"front-office/pkg/core/template"
	"front-office/pkg/procat"
	"front-office/pkg/procat/phonelivestatus"
	"front-office/pkg/scoreezy/genretail"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(routeGroup fiber.Router, cfg *config.Config) {
	userGroup := routeGroup.Group("users")
	auth.SetupInit(userGroup, cfg)
	member.SetupInit(userGroup, cfg)

	roleGroup := routeGroup.Group("roles")
	role.SetupInit(roleGroup, cfg)

	permissionGroup := routeGroup.Group("permissions")
	permission.SetupInit(permissionGroup)

	companyGroup := routeGroup.Group("companies")
	company.SetupInit(companyGroup)

	gradingGroup := routeGroup.Group("gradings")
	grading.SetupInit(gradingGroup, cfg)

	genRetailGroup := routeGroup.Group("scores")
	genretail.SetupInit(genRetailGroup, cfg)

	logGroup := routeGroup.Group("logs")
	transaction.SetupInit(logGroup, cfg)
	operation.SetupInit(logGroup, cfg)

	productGroup := routeGroup.Group("products")
	procat.SetupInit(productGroup, cfg)
	phonelivestatus.SetupInit(productGroup, cfg)

	templateGroup := routeGroup.Group("templates")
	template.SetupInit(templateGroup)
}

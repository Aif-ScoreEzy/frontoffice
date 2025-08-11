package core

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/auth"
	"front-office/pkg/core/grade"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"
	"front-office/pkg/core/template"
	"front-office/pkg/procat"
	"front-office/pkg/scoreezy/genretail"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(routeGroup fiber.Router, cfg *config.Config) {
	client := httpclient.NewDefaultClient(10 * time.Second)

	userGroup := routeGroup.Group("users")
	auth.SetupInit(userGroup, cfg, client)
	member.SetupInit(userGroup, cfg, client)

	roleGroup := routeGroup.Group("roles")
	role.SetupInit(roleGroup, cfg, client)

	gradeGroup := routeGroup.Group("grades")
	grade.SetupInit(gradeGroup, cfg, client)

	genRetailGroup := routeGroup.Group("scoreezy")
	genretail.SetupInit(genRetailGroup, cfg, client)

	logGroup := routeGroup.Group("logs")
	transaction.SetupInit(logGroup, cfg, client)
	operation.SetupInit(logGroup, cfg, client)

	productGroup := routeGroup.Group("products")
	procat.SetupInit(productGroup, cfg)

	templateGroup := routeGroup.Group("templates")
	template.SetupInit(templateGroup)
}

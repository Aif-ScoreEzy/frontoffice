package transaction

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(logAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	roleRepository := role.NewRepository(cfg)
	roleService := role.NewService(roleRepository)

	memberRepository := member.NewRepository(cfg, client)
	memberService := member.NewService(memberRepository, roleService)

	repository := NewRepository(cfg, client)
	service := NewService(repository)
	controller := NewController(service, memberService)

	logTransScoreezyAPI := logAPI.Group("scoreezy")
	logTransScoreezyAPI.Get("/", controller.GetLogScoreezy)
	logTransScoreezyAPI.Get("/by-date", controller.GetLogScoreezyByDate)
	logTransScoreezyAPI.Get("/by-range-date", controller.GetLogScoreezyByRangeDate)
	logTransScoreezyAPI.Get("/by-month", controller.GetLogScoreezyByMonth)

	// logTransProcatGroup := logAPI.Group("product_catalog")
}

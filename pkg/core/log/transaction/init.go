package transaction

import (
	"front-office/app/config"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(logAPI fiber.Router, cfg *config.Config) {
	roleRepository := role.NewRepository(cfg)
	roleService := role.NewService(roleRepository)

	memberRepository := member.NewRepository(cfg)
	memberService := member.NewService(memberRepository, roleService)

	repository := NewRepository(cfg)
	service := NewService(repository, cfg)
	controller := NewController(service, memberService)

	logTransScoreezyAPI := logAPI.Group("scoreezy")
	logTransScoreezyAPI.Get("/", controller.GetLogTransactions)
	logTransScoreezyAPI.Get("/by-date", controller.GetLogTransactionsByDate)
	logTransScoreezyAPI.Get("/by-range-date", controller.GetLogTransactionsByRangeDate)
	logTransScoreezyAPI.Get("/by-month", controller.GetLogTransactionsByMonth)
}

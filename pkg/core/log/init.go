package log

import (
	"front-office/app/config"
	"front-office/pkg/core/member"
	"front-office/pkg/core/role"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(logAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	roleRepository := role.NewRepository(db, cfg)
	roleService := role.NewService(roleRepository)

	memberRepository := member.NewRepository(db, cfg)
	memberService := member.NewService(memberRepository, roleService)

	repository := NewRepository(db, cfg)
	service := NewService(repository, cfg)
	controller := NewController(service, memberService)

	// scoreezy
	logTransScoreezyAPI := logAPI.Group("scoreezy")
	logTransScoreezyAPI.Get("/", controller.GetLogTransactions)
	logTransScoreezyAPI.Get("/by-date", controller.GetLogTransactionsByDate)
	logTransScoreezyAPI.Get("/by-range-date", controller.GetLogTransactionsByRangeDate)
	logTransScoreezyAPI.Get("/by-month", controller.GetLogTransactionsByMonth)
}

package procat

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/procat/loanrecordchecker"
	"front-office/pkg/procat/log"
	"front-office/pkg/procat/multipleloan"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(routeAPI fiber.Router, cfg *config.Config) {
	client := httpclient.NewDefaultClient(10 * time.Second)

	complianceGroupAPI := routeAPI.Group("compliance")
	loanrecordchecker.SetupInit(complianceGroupAPI, cfg, client)
	multipleloan.SetupInit(complianceGroupAPI, cfg, client)
	log.SetupInit(complianceGroupAPI, cfg, client)
}

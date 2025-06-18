package procat

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/procat/compliance/loanrecordchecker"
	"front-office/pkg/procat/compliance/multipleloan"
	taxcompliancestatus "front-office/pkg/procat/incometax/taxcompliencestatus"
	"front-office/pkg/procat/incometax/taxpayerstatus"
	"front-office/pkg/procat/incometax/taxscore"
	"front-office/pkg/procat/incometax/taxverificationdetail"
	"front-office/pkg/procat/log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(routeAPI fiber.Router, cfg *config.Config) {
	client := httpclient.NewDefaultClient(10 * time.Second)

	complianceGroupAPI := routeAPI.Group("compliance")
	loanrecordchecker.SetupInit(complianceGroupAPI, cfg, client)
	multipleloan.SetupInit(complianceGroupAPI, cfg, client)
	log.SetupInit(complianceGroupAPI, cfg, client)

	incomeTaxGroupAPI := routeAPI.Group("incometax")
	taxcompliancestatus.SetupInit(incomeTaxGroupAPI, cfg, client)
	taxpayerstatus.SetupInit(incomeTaxGroupAPI, cfg, client)
	taxscore.SetupInit(incomeTaxGroupAPI, cfg, client)
	taxverificationdetail.SetupInit(incomeTaxGroupAPI, cfg, client)
	log.SetupInit(incomeTaxGroupAPI, cfg, client)
}

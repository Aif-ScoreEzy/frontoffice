package job

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client)
	transactionRepo := transaction.NewRepository(cfg, client)
	service := NewService(repository, transactionRepo)
	controller := NewController(service)

	apiGroup.Get("/:product_slug/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetProCatJob)
	apiGroup.Get("/:product_slug/jobs/:job_id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetProCatJobDetail)
}

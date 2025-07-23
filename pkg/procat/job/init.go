package job

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client, nil)
	transactionRepo := transaction.NewRepository(cfg, client, nil)
	service := NewService(repository, transactionRepo)
	controller := NewController(service)

	apiGroup.Get("/:product_slug/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJob)
	apiGroup.Get("/:product_slug/jobs/:job_id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	apiGroup.Get("/:product_slug/jobs/:job_id/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetails)
	apiGroup.Get("/:product_slug/job-details", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetailsByDateRange)
	apiGroup.Get("/:product_slug/job-details/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetailsByDateRange)
}

package taxcompliancestatus

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/middleware"
	"front-office/pkg/procat/job"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	productRepo := product.NewRepository(cfg, client)
	jobRepo := job.NewRepository(cfg, client, nil)
	transactionRepo := transaction.NewRepository(cfg, client, nil)

	jobService := job.NewService(jobRepo, transactionRepo)
	service := NewService(repo, productRepo, jobRepo, transactionRepo, jobService)

	controller := NewController(service)

	taxComplianceGroup := apiGroup.Group("tax-compliance-status")
	taxComplianceGroup.Post("/", middleware.Auth(), middleware.IsRequestValid(taxComplianceStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.TaxComplianceStatus)
}

package multipleloan

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
	transactionRepo := transaction.NewRepository(cfg, client)

	jobService := job.NewService(jobRepo, transactionRepo)
	service := NewService(repo, productRepo, jobRepo, transactionRepo, jobService)

	controller := NewController(service)

	apiGroup.Post("/:product_slug/single-request", middleware.Auth(), middleware.IsRequestValid(multipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan)
	apiGroup.Post("/:product_slug/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkMultipleLoan)
}

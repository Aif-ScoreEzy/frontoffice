package phonelivestatus

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
	repository := NewRepository(cfg, client, nil)
	productRepo := product.NewRepository(cfg, client)
	jobRepo := job.NewRepository(cfg, client, nil)
	transactionRepo := transaction.NewRepository(cfg, client, nil)
	jobService := job.NewService(jobRepo, transactionRepo)
	service := NewService(repository, productRepo, jobRepo, transactionRepo, jobService)
	controller := NewController(service)

	phoneLiveStatusGroup := apiGroup.Group("phone-live-status")
	phoneLiveStatusGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(phoneLiveStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.SingleSearch)
	phoneLiveStatusGroup.Post("/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkSearch)
	phoneLiveStatusGroup.Get("/jobs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobs)
	phoneLiveStatusGroup.Get("/jobs/:id/details", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobDetails)
	phoneLiveStatusGroup.Get("/jobs/:id/details/export", middleware.AdminAuth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobDetails)
	phoneLiveStatusGroup.Get("/jobs-summary", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetJobsSummary)
	phoneLiveStatusGroup.Get("/jobs-summary/export", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.ExportJobsSummary)
}

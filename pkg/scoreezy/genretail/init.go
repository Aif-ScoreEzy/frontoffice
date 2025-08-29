package genretail

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/grade"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	gradeRepo := grade.NewRepository(cfg, client, nil)
	transRepo := transaction.NewRepository(cfg, client, nil)
	productRepo := product.NewRepository(cfg, client)
	logRepo := operation.NewRepository(cfg, client, nil)

	service := NewService(repo, gradeRepo, transRepo, productRepo, logRepo)

	controller := NewController(service)

	genRetailGroup := apiGroup.Group("gen-retail")
	genRetailGroup.Post("/dummy-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(genRetailRequest{}), controller.DummyRequestScore)
	genRetailGroup.Post("/single-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(genRetailRequest{}), controller.SingleRequest)
	genRetailGroup.Post("/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkRequest)
	genRetailGroup.Get("/logs", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetLogsScoreezy)
	genRetailGroup.Get("/logs/:trx_id", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetLogScoreezy)
	// genRetailAPI.Put("/upload-scoring-template", middleware.Auth(), middleware.IsRequestValid(UploadScoringRequest{}), middleware.GetJWTPayloadFromCookie(), middleware.DocUpload(), controller.UploadCSV)
}

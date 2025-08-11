package genretail

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/grade"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	gradeRepo := grade.NewRepository(cfg, client, nil)
	logOpRepo := operation.NewRepository(cfg, client, nil)

	service := NewService(repo, gradeRepo)
	logOpService := operation.NewService(logOpRepo)

	controller := NewController(service, logOpService)

	genRetailGroup := apiGroup.Group("gen-retail")
	genRetailGroup.Post("/single-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(genRetailRequest{}), controller.RequestScore)
	// genRetailAPI.Get("/download-csv/:opsi", middleware.GetJWTPayloadFromCookie(), controller.DownloadCSV)
	// genRetailAPI.Put("/upload-scoring-template", middleware.Auth(), middleware.IsRequestValid(UploadScoringRequest{}), middleware.GetJWTPayloadFromCookie(), middleware.DocUpload(), controller.UploadCSV)
	// genRetailAPI.Get("/bulk-search", middleware.GetJWTPayloadFromCookie(), controller.GetBulkSearch)
}

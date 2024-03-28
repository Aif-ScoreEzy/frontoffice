package genretail

import (
	"front-office/app/config"
	"front-office/pkg/core/grading"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(genRetailAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	repoGrading := grading.NewRepository(db)
	service := NewService(repo, cfg)
	serviceGrading := grading.NewService(repoGrading)
	controller := NewController(service, serviceGrading)

	genRetailAPI.Post("/request-score", middleware.GetPayloadFromJWT(), middleware.IsRequestValid(GenRetailRequest{}), controller.RequestScore)
	genRetailAPI.Get("/download-csv/:opsi", middleware.GetPayloadFromJWT(), controller.DownloadCSV)
	genRetailAPI.Put("/upload-scoring-template", middleware.Auth(), middleware.IsRequestValid(UploadScoringRequest{}), middleware.GetPayloadFromJWT(), middleware.DocUpload(), controller.UploadCSV)
	genRetailAPI.Get("/bulk-search", middleware.GetPayloadFromJWT(), controller.GetBulkSearch)
}

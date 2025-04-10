package genretail

import (
	"front-office/app/config"
	"front-office/pkg/core/grading"
	"front-office/pkg/core/log/operation"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(genRetailAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	repoGrading := grading.NewRepository(db, cfg)
	repoLogOperation := operation.NewRepository(cfg)

	service := NewService(repo, cfg)
	serviceGrading := grading.NewService(repoGrading)
	serviceLogOperation := operation.NewService(repoLogOperation)

	controller := NewController(service, serviceGrading, serviceLogOperation)

	genRetailAPI.Post("/request-score", middleware.GetJWTPayloadFromCookie(), middleware.IsRequestValid(GenRetailRequest{}), controller.RequestScore)
	genRetailAPI.Get("/download-csv/:opsi", middleware.GetJWTPayloadFromCookie(), controller.DownloadCSV)
	genRetailAPI.Put("/upload-scoring-template", middleware.Auth(), middleware.IsRequestValid(UploadScoringRequest{}), middleware.GetJWTPayloadFromCookie(), middleware.DocUpload(), controller.UploadCSV)
	genRetailAPI.Get("/bulk-search", middleware.GetJWTPayloadFromCookie(), controller.GetBulkSearch)
}

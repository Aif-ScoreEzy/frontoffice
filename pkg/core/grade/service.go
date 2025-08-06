package grade

import (
	"fmt"
	"front-office/internal/apperror"
)

func NewService(repo Repository) Service {
	return &service{Repo: repo}
}

type service struct {
	Repo Repository
}

type Service interface {
	CreateGrading(payload *createGradePayload) error
}

func (svc *service) CreateGrading(payload *createGradePayload) error {
	for i := 0; i < len(payload.Request.Grades); i++ {
		for j := i + 1; j < len(payload.Request.Grades); j++ {
			if payload.Request.Grades[i].Grade == payload.Request.Grades[j].Grade {
				return apperror.BadRequest(fmt.Sprintf("duplicate grade: %s", payload.Request.Grades[i].Grade))
			}
			if !(payload.Request.Grades[i].End <= payload.Request.Grades[j].Start || payload.Request.Grades[j].End <= payload.Request.Grades[i].Start) {
				return apperror.BadRequest(fmt.Sprintf("overlapping grade range between %s and %s", payload.Request.Grades[i].Grade, payload.Request.Grades[j].Grade))
			}
		}
	}

	if err := svc.Repo.CreateGradesAPI(payload); err != nil {
		return apperror.MapRepoError(err, "failed to create grades")
	}

	return nil
}

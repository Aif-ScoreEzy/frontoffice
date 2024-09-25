package member

import (
	"encoding/json"
	"io"
)

func NewService(repo Repository) Service {
	return &service{
		Repo: repo,
	}
}

type service struct {
	Repo Repository
}

type Service interface {
	GetBy(query *FindUserQuery) (*AifResponse, error)
}

func (s *service) GetBy(query *FindUserQuery) (*AifResponse, error) {
	response, err := s.Repo.GetBy(query)
	if err != nil {
		return nil, err
	}

	var baseResponse *AifResponse
	if response != nil {
		dataByte, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		if err := json.Unmarshal(dataByte, &baseResponse); err != nil {
			return nil, err
		}
	}

	return baseResponse, nil
}

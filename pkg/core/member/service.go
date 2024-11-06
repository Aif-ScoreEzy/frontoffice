package member

import (
	"encoding/json"
	"io"
	"net/http"
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
	GetMemberBy(query *FindUserQuery) (*AifResponse, error)
	GetMemberList() (*AifResponse, error)
}

func (s *service) GetMemberBy(query *FindUserQuery) (*AifResponse, error) {
	response, err := s.Repo.GetMemberBy(query)
	if err != nil {
		return nil, err
	}

	return s.parseResponse(response)
}

func (s *service) GetMemberList() (*AifResponse, error) {
	response, err := s.Repo.GetMemberList()
	if err != nil {
		return nil, err
	}

	return s.parseResponse(response)
}

func (s *service) parseResponse(response *http.Response) (*AifResponse, error) {
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

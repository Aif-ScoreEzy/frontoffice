package multipleloan

import (
	"encoding/json"
	"front-office/app/config"
	"io"
	"net/http"
)

func NewService(cfg *config.Config, repo Repository) Service {
	return &service{
		Cfg:  cfg,
		Repo: repo,
	}
}

type service struct {
	Cfg  *config.Config
	Repo Repository
}

type Service interface {
	CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*MultipleLoanRawResponse, error)
}

func (svc *service) CallMultipleLoan7Days(request *MultipleLoanRequest, apiKey string) (*MultipleLoanRawResponse, error) {
	response, err := svc.Repo.CallMultipleLoan7Days(request, apiKey)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseResponse(response *http.Response) (*MultipleLoanRawResponse, error) {
	var baseResponse *MultipleLoanRawResponse

	if response != nil {
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		if err := json.Unmarshal(dataBytes, &baseResponse); err != nil {
			return nil, err
		}
	}

	return baseResponse, nil
}

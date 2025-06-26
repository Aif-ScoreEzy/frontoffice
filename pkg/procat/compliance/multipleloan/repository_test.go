package multipleloan

import (
	"bytes"
	"encoding/json"
	"errors"
	"front-office/app/config"
	"front-office/common/constant"
	"front-office/common/model"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCallMultipleLoan7DaysAPI_Success(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}, mockClient)

	mockData := model.ProCatAPIResponse[dataMultipleLoanResponse]{
		Success: true,
		Message: "ok",
		Data:    dataMultipleLoanResponse{},
	}
	body, _ := json.Marshal(mockData)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan7Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan7DaysAPI_NewRequestError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
	}, mockClient)

	_, err := repo.CallMultipleLoan7Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Error(t, err)
}

func TestCallMultipleLoan7DaysAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	req := &multipleLoanRequest{}

	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallMultipleLoan7Days("test-api-key", "job-id", "member-id", "company-id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan7DaysAPI_ParseError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan7Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan30DaysAPI_Success(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}, mockClient)

	mockData := model.ProCatAPIResponse[dataMultipleLoanResponse]{
		Success: true,
		Message: "ok",
		Data:    dataMultipleLoanResponse{},
	}
	body, _ := json.Marshal(mockData)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan30Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan30DaysAPI_NewRequestError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
	}, mockClient)

	_, err := repo.CallMultipleLoan30Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Error(t, err)
}

func TestCallMultipleLoan30DaysAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	req := &multipleLoanRequest{}

	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallMultipleLoan30Days("test-api-key", "job-id", "member-id", "company-id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan30DaysAPI_ParseError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan30Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan90DaysAPI_Success(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}, mockClient)

	mockData := model.ProCatAPIResponse[dataMultipleLoanResponse]{
		Success: true,
		Message: "ok",
		Data:    dataMultipleLoanResponse{},
	}
	body, _ := json.Marshal(mockData)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan90Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan90DaysAPI_NewRequestError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
	}, mockClient)

	_, err := repo.CallMultipleLoan90Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Error(t, err)
}

func TestCallMultipleLoan90DaysAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	req := &multipleLoanRequest{}

	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallMultipleLoan90Days("test-api-key", "job-id", "member-id", "company-id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestCallMultipleLoan90DaysAPI_ParseError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallMultipleLoan90Days("apiKey", "jobId", "memberId", "companyId", &multipleLoanRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

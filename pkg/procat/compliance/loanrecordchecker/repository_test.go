package loanrecordchecker

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

func TestCallLoanRecordCheckerAPI_Success(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}, mockClient)

	mockData := model.ProCatAPIResponse[dataLoanRecord]{
		Success: true,
		Message: "ok",
		Data:    dataLoanRecord{},
	}
	body, _ := json.Marshal(mockData)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallLoanRecordCheckerAPI("apiKey", "jobId", "memberId", "companyId", &loanRecordCheckerRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	mockClient.AssertExpectations(t)
}

func TestCallLoanRecordCheckerAPI_NewRequestError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
	}, mockClient)

	_, err := repo.CallLoanRecordCheckerAPI("apiKey", "jobId", "memberId", "companyId", &loanRecordCheckerRequest{})
	assert.Error(t, err)
}

func TestCallPhoneLiveStatusAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	req := &loanRecordCheckerRequest{}

	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get(constant.XAPIKey) == "test-api-key" &&
				req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallLoanRecordCheckerAPI("test-api-key", "job-id", "member-id", "company-id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestCallLoanRecordCheckerAPI_ParseError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallLoanRecordCheckerAPI("apiKey", "jobId", "memberId", "companyId", &loanRecordCheckerRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

package taxscore

import (
	"bytes"
	"errors"
	"front-office/app/config"
	"front-office/common/constant"
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

func TestCallTaxScoreAPI_Success(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	expectedBody := `{
		"message": "Succeed to Request Data.",
		"success": true,
		"input": {
			"npwp": "0658552450011000"
		},
		"data": {
			"nama": "",
			"alamat": "",
			"score": "0",
			"status": ""
		},
		"pricing_strategy": "PAY",
		"transaction_id": "c2a97f8d-b482-4c5b-98fb-cee86a0e122e",
		"datetime": "2025-06-26 17:04:21"
	}`

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(expectedBody))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(mockResp, nil)

	req := &taxScoreRequest{
		Npwp: "092542823407000",
	}

	resp, err := repo.CallTaxScoreAPI("test-api-key", "1", req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestCallTaxScoreAPI_NewRequestError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
	}, mockClient)

	_, err := repo.CallTaxScoreAPI("apiKey", "jobId", &taxScoreRequest{})
	assert.Error(t, err)
}

func TestCallTaxScoreAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallTaxScoreAPI("test-api-key", "job-id", &taxScoreRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestCallTaxScoreAPI_ParseError(t *testing.T) {
	mockClient := new(MockClient)
	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
	}

	mockClient.On("Do", mock.Anything).Return(resp, nil)

	result, err := repo.CallTaxScoreAPI("apiKey", "jobId", &taxScoreRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

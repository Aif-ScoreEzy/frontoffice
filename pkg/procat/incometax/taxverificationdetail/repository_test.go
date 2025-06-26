package taxverificationdetail

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
			"npwp": "0092542823407000"
		},
		"data": {
			"nama": "",
			"alamat": "",
			"status": "Unreported"
		},
		"pricing_strategy": "PAY",
		"transaction_id": "9c6b46c9-e3be-4c90-a5e3-894b26432e0b",
		"datetime": "2025-06-01 05:17:53"
	}`

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(expectedBody))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(mockResp, nil)

	req := &taxVerificationRequest{
		NpwpOrNik: "092542823407000",
	}

	resp, err := repo.CallTaxVerificationAPI("test-api-key", "1", req)

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

	_, err := repo.CallTaxVerificationAPI("apiKey", "jobId", &taxVerificationRequest{})
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

	_, err := repo.CallTaxVerificationAPI("test-api-key", "job-id", &taxVerificationRequest{})

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

	result, err := repo.CallTaxVerificationAPI("apiKey", "jobId", &taxVerificationRequest{})
	assert.Nil(t, result)
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}

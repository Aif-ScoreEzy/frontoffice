package phonelivestatus

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

func TestPhoneLiveStatusAPI_Success(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	expectedBody := `{
			"message": "Succeed to Request Data",
			"success": true,
			"input": {
				"phone number": "081282177383"
			},
			"data": {
				"live_status": "active, reachable"
			},
			"pricing_strategy": "PAY",
			"transaction_id": "",
			"date": "yyyy-mmm-dd HH-MM-SS"
		}`

	mockResp := &http.Response{
		Body:   io.NopCloser(bytes.NewReader([]byte(expectedBody))),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(mockResp, nil)

	req := &phoneLiveStatusRequest{
		PhoneNumber: "08111111110",
		TrxId:       "trx-123",
	}

	resp, err := repo.CallPhoneLiveStatusAPI("test-api-key", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockClient.AssertExpectations(t)
}

func TestCallPhoneLiveStatusAPI_HTTPRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	// Setup request yang valid
	req := &phoneLiveStatusRequest{}

	// Setup mock untuk mengembalikan error
	expectedErr := errors.New("failed to make HTTP request")
	mockClient.
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Header.Get(constant.XAPIKey) == "test-api-key" &&
				req.Header.Get("Content-Type") == "application/json"
		})).
		Return(&http.Response{}, expectedErr)

	_, err := repo.CallPhoneLiveStatusAPI("test-api-key", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make HTTP request")
	mockClient.AssertExpectations(t)
}

func TestPhoneLiveStatusAPI_NewRequestError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockInvalidHost, // Invalid URL to simulate error
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	req := &phoneLiveStatusRequest{
		PhoneNumber: "08111111110",
		TrxId:       "trx-123",
	}

	_, err := repo.CallPhoneLiveStatusAPI("test-api-key", req)

	assert.Error(t, err)
}

func TestPhoneLiveStatusAPI_ParseResponseError(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	invalidJSON := `this is not json`

	mockResp := &http.Response{
		Body:       io.NopCloser(strings.NewReader(invalidJSON)),
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(mockResp, nil)

	req := &phoneLiveStatusRequest{
		PhoneNumber: "08111111110",
		TrxId:       "trx-123",
	}

	_, err := repo.CallPhoneLiveStatusAPI("test-api-key", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

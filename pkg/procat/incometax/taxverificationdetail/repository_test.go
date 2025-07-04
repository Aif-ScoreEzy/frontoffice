package taxverificationdetail

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
	"github.com/stretchr/testify/require"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func setupMockRepo(t *testing.T, response *http.Response, err error) (Repository, *MockClient) {
	t.Helper()

	mockClient := new(MockClient)
	mockClient.On("Do", mock.Anything).Return(response, err)

	repo := NewRepository(&config.Config{
		Env: &config.Environment{ProductCatalogHost: constant.MockHost},
	}, mockClient, nil)

	return repo, mockClient
}

func TestCallTaxVerificationAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.ProCatAPIResponse[taxVerificationRespData]{
			Success: true,
			Message: "Succeed to Request Data.",
			Data: taxVerificationRespData{
				Status: "Unreported",
			},
			PricingStrategy: "PAY",
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallTaxVerificationAPI(constant.DummyAPIKey, constant.DummyJobId, &taxVerificationRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, "Succeed to Request Data.", result.Message)
		assert.Equal(t, "Unreported", result.Data.Status)
		assert.Equal(t, "PAY", result.PricingStrategy)
		mockClient.AssertExpectations(t)
	})

	t.Run("MarshalError", func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{ProductCatalogHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		result, err := repo.CallTaxVerificationAPI(constant.DummyAPIKey, constant.DummyJobId, &taxVerificationRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallTaxVerificationAPI(constant.DummyAPIKey, constant.DummyJobId, &taxVerificationRequest{})
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		req := &taxVerificationRequest{}
		_, err := repo.CallTaxVerificationAPI(constant.DummyAPIKey, constant.DummyJobId, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run("ParseError", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallTaxVerificationAPI(constant.DummyAPIKey, constant.DummyJobId, &taxVerificationRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

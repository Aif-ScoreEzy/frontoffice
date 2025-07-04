package taxscore

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

func TestCallTaxScoreAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.ProCatAPIResponse[taxScoreRespData]{
			Success: true,
			Message: "Succeed to Request Data.",
			Data: taxScoreRespData{
				Score: "0",
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

		result, err := repo.CallTaxScoreAPI(constant.DummyAPIKey, constant.DummyJobId, &taxScoreRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, "Succeed to Request Data.", result.Message)
		assert.Equal(t, "0", result.Data.Score)
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

		result, err := repo.CallTaxScoreAPI(constant.DummyAPIKey, constant.DummyJobId, &taxScoreRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallTaxScoreAPI(constant.DummyAPIKey, constant.DummyJobId, &taxScoreRequest{})
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		req := &taxScoreRequest{}
		_, err := repo.CallTaxScoreAPI(constant.DummyAPIKey, constant.DummyJobId, req)

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

		result, err := repo.CallTaxScoreAPI(constant.DummyAPIKey, constant.DummyJobId, &taxScoreRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

package transaction

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
		Env: &config.Environment{AifcoreHost: constant.MockHost},
	}, mockClient, nil)

	return repo, mockClient
}

func TestGetLogTransByJobIdAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[[]*LogTransProductCatalog]{
			Success: true,
			Data:    []*LogTransProductCatalog{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.GetLogTransByJobIdAPI(constant.DummyJobId, constant.DummyCompanyId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		result, err := repo.GetLogTransByJobIdAPI(constant.DummyJobId, constant.DummyCompanyId)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		result, err := repo.GetLogTransByJobIdAPI(constant.DummyJobId, constant.DummyCompanyId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run("ParseError", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.GetLogTransByJobIdAPI(constant.DummyJobId, constant.DummyCompanyId)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestProcessedLogCountAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[*getProcessedCountResp]{
			Success: true,
			Data: &getProcessedCountResp{
				ProcessedCount: uint(constant.DummyCount),
			},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.ProcessedLogCountAPI(constant.DummyJobId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(constant.DummyCount), uint(result.ProcessedCount))
		mockClient.AssertExpectations(t)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		result, err := repo.ProcessedLogCountAPI(constant.DummyJobId)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		result, err := repo.ProcessedLogCountAPI(constant.DummyJobId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run("ParseError", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.ProcessedLogCountAPI(constant.DummyJobId)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCreateLogTransAPI(t *testing.T) {
	addLogReq := &LogTransProCatRequest{}

	t.Run("Success", func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[any]{
			Success: true,
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		err = repo.CreateLogTransAPI(addLogReq)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("MarshalError", func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		err := repo.CreateLogTransAPI(addLogReq)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.CreateLogTransAPI(addLogReq)

		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.CreateLogTransAPI(addLogReq)

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

		err := repo.CreateLogTransAPI(addLogReq)

		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

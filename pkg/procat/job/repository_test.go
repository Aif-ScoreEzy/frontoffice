package job

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

func TestCallCreateProCatJob(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[any]{
			Success: true,
			Message: "Succeed to Request Data.",
			Data: createJobRespData{
				JobId: 1,
			},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallCreateJobAPI(&CreateJobRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result.JobId, uint(1))
		mockClient.AssertExpectations(t)
	})

	t.Run("MarshalError", func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		result, err := repo.CallCreateJobAPI(&CreateJobRequest{})

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallCreateJobAPI(&CreateJobRequest{})
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		req := &CreateJobRequest{}
		_, err := repo.CallCreateJobAPI(req)

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

		result, err := repo.CallCreateJobAPI(&CreateJobRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallUpdateJob(t *testing.T) {
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

		err = repo.CallUpdateJob(constant.DummyJobId, map[string]interface{}{})

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

		err := repo.CallUpdateJob(constant.DummyJobId, map[string]interface{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.CallUpdateJob(constant.DummyJobId, map[string]interface{}{})

		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.CallUpdateJob(constant.DummyJobId, map[string]interface{}{})

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

		err := repo.CallUpdateJob(constant.DummyJobId, map[string]interface{}{})

		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetProCatJobAPI(t *testing.T) {
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

		result, err := repo.CallGetJobsAPI(&logFilter{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetJobsAPI(&logFilter{})

		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetJobsAPI(&logFilter{})

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

		result, err := repo.CallGetJobsAPI(&logFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetProCatJobDetailAPI(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[any]{
			Success: true,
			Data: &jobDetailResponse{
				TotalData: 3,
			},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetJobDetailAPI(&logFilter{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result.Data.TotalData, int64(3))
		mockClient.AssertExpectations(t)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetJobDetailAPI(&logFilter{})

		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetJobDetailAPI(&logFilter{})

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

		result, err := repo.CallGetJobDetailAPI(&logFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

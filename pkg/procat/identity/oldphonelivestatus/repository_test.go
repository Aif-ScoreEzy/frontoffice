package oldphonelivestatus

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
		Env: &config.Environment{
			ProductCatalogHost: constant.MockHost,
			AifcoreHost:        constant.MockHost,
		},
	}, mockClient, nil)

	return repo, mockClient
}

func TestCallCreateJobAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[createJobRespData]{
			Success: true,
			Data:    createJobRespData{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CreateJobAPI(&createJobRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseMarshalError, func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		result, err := repo.CreateJobAPI(&createJobRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CreateJobAPI(&createJobRequest{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CreateJobAPI(&createJobRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CreateJobAPI(&createJobRequest{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetPhoneLiveStatusJobAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[jobListRespData]{
			Success: true,
			Data:    jobListRespData{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			JobId: "100",
		}
		result, err := repo.CallGetPhoneLiveStatusJobAPI(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "100", filter.JobId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetPhoneLiveStatusJobAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetPhoneLiveStatusJobAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetPhoneLiveStatusJobAPI(&phoneLiveStatusFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetJobDetailsAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[jobDetailRespData]{
			Success: true,
			Data:    jobDetailRespData{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			JobId: "100",
		}
		result, err := repo.GetJobDetailsAPI(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "100", filter.JobId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.GetJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.GetJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.GetJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetAllJobDetailsAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[[]*mstPhoneLiveStatusJobDetail]{
			Success: true,
			Data:    []*mstPhoneLiveStatusJobDetail{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			JobId: "100",
		}
		result, err := repo.CallGetAllJobDetailsAPI(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "100", filter.JobId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetAllJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetAllJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetAllJobDetailsAPI(&phoneLiveStatusFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetJobDetailsByDateRangeAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[[]*mstPhoneLiveStatusJobDetail]{
			Success: true,
			Data:    []*mstPhoneLiveStatusJobDetail{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			MemberId: constant.DummyMemberId,
		}
		result, err := repo.CallGetJobDetailsByDateRangeAPI(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, constant.DummyMemberId, filter.MemberId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetJobDetailsByDateRangeAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetJobDetailsByDateRangeAPI(&phoneLiveStatusFilter{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetJobDetailsByDateRangeAPI(&phoneLiveStatusFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetJobsSummary(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[jobsSummaryRespData]{
			Success: true,
			Data:    jobsSummaryRespData{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			MemberId: constant.DummyMemberId,
		}
		result, err := repo.CallGetJobsSummary(filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, constant.DummyMemberId, filter.MemberId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetJobsSummary(&phoneLiveStatusFilter{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetJobsSummary(&phoneLiveStatusFilter{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetJobsSummary(&phoneLiveStatusFilter{})

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallGetProcessedCount(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[jobsSummaryRespData]{
			Success: true,
			Data:    jobsSummaryRespData{},
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		filter := &phoneLiveStatusFilter{
			MemberId: constant.DummyMemberId,
		}
		result, err := repo.CallGetProcessedCountAPI(constant.DummyJobId)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, constant.DummyMemberId, filter.MemberId)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallGetProcessedCountAPI(constant.DummyJobId)

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		_, err := repo.CallGetProcessedCountAPI(constant.DummyJobId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallGetProcessedCountAPI(constant.DummyJobId)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallUpdateJob(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[any]{
			Success: true,
			Data:    nil,
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		successStatus := "success"
		req := &updateJobRequest{
			Status: &successStatus,
		}
		err = repo.UpdateJobAPI(constant.DummyJobId, req)

		assert.NoError(t, err)
		assert.Equal(t, &successStatus, req.Status)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseMarshalError, func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		err := repo.UpdateJobAPI(constant.DummyJobId, &updateJobRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.UpdateJobAPI(constant.DummyJobId, &updateJobRequest{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.UpdateJobAPI(constant.DummyJobId, &updateJobRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		err := repo.UpdateJobAPI(constant.DummyJobId, &updateJobRequest{})

		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallUpdateJobDetail(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.AifcoreAPIResponse[any]{
			Success: true,
			Data:    nil,
		}
		body, err := json.Marshal(mockData)
		require.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		successStatus := "success"
		req := &updateJobDetailRequest{
			Status: &successStatus,
		}
		err = repo.CallUpdateJobDetail(constant.DummyJobId, constant.DummyJobDetailId, req)

		assert.NoError(t, err)
		assert.Equal(t, &successStatus, req.Status)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseMarshalError, func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		err := repo.CallUpdateJobDetail(constant.DummyJobId, constant.DummyJobDetailId, &updateJobDetailRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.CallUpdateJobDetail(constant.DummyJobId, constant.DummyJobDetailId, &updateJobDetailRequest{})

		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.CallUpdateJobDetail(constant.DummyJobId, constant.DummyJobDetailId, &updateJobDetailRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		err := repo.CallUpdateJobDetail(constant.DummyJobId, constant.DummyJobDetailId, &updateJobDetailRequest{})

		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestCallPhoneLiveStatusAPI(t *testing.T) {
	t.Run(constant.TestCaseSuccess, func(t *testing.T) {
		mockData := model.ProCatAPIResponse[phoneLiveStatusRespData]{
			Success: true,
			Message: "Succeed to Request Data.",
			Data: phoneLiveStatusRespData{
				LiveStatus: "active, reachable",
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

		result, err := repo.CallPhoneLiveStatusAPI(constant.DummyAPIKey, &phoneLiveStatusRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, "Succeed to Request Data.", result.Message)
		assert.Equal(t, "active, reachable", result.Data.LiveStatus)
		assert.Equal(t, "PAY", result.PricingStrategy)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseMarshalError, func(t *testing.T) {
		fakeMarshal := func(v any) ([]byte, error) {
			return nil, errors.New(constant.ErrFailedMarshalReq)
		}

		repo := NewRepository(&config.Config{
			Env: &config.Environment{ProductCatalogHost: constant.MockHost},
		}, &MockClient{}, fakeMarshal)

		result, err := repo.CallPhoneLiveStatusAPI(constant.DummyAPIKey, &phoneLiveStatusRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run(constant.TestCaseNewRequestError, func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{ProductCatalogHost: constant.MockInvalidHost},
		}, mockClient, nil)

		_, err := repo.CallPhoneLiveStatusAPI(constant.DummyAPIKey, &phoneLiveStatusRequest{})
		assert.Error(t, err)
	})

	t.Run(constant.TestCaseHTTPRequestError, func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		req := &phoneLiveStatusRequest{}
		_, err := repo.CallPhoneLiveStatusAPI(constant.DummyAPIKey, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrHTTPReqFailed)
		mockClient.AssertExpectations(t)
	})

	t.Run(constant.TestCaseParseError, func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{invalid-json`)),
		}

		repo, mockClient := setupMockRepo(t, resp, nil)

		result, err := repo.CallPhoneLiveStatusAPI(constant.DummyAPIKey, &phoneLiveStatusRequest{})
		assert.Nil(t, result)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

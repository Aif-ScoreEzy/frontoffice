package auth

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

func TestVerifyMemberAPI(t *testing.T) {
	passwordResetReq := &PasswordResetRequest{
		Password:        constant.DummyPassword,
		ConfirmPassword: constant.DummyPassword,
	}

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

		err = repo.VerifyMemberAPI(constant.DummyMemberId, passwordResetReq)

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

		err := repo.VerifyMemberAPI(constant.DummyMemberId, passwordResetReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.VerifyMemberAPI(constant.DummyMemberId, passwordResetReq)
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.VerifyMemberAPI(constant.DummyMemberId, passwordResetReq)

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

		err := repo.VerifyMemberAPI(constant.DummyMemberId, passwordResetReq)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestPasswordResetAPI(t *testing.T) {
	passwordResetReq := &PasswordResetRequest{
		Password:        constant.DummyPassword,
		ConfirmPassword: constant.DummyPassword,
	}

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

		err = repo.PasswordResetAPI(constant.DummyMemberId, constant.DummyToken, passwordResetReq)

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

		err := repo.PasswordResetAPI(constant.DummyMemberId, constant.DummyToken, passwordResetReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.PasswordResetAPI(constant.DummyMemberId, constant.DummyToken, passwordResetReq)
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.PasswordResetAPI(constant.DummyMemberId, constant.DummyToken, passwordResetReq)

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

		err := repo.PasswordResetAPI(constant.DummyMemberId, constant.DummyToken, passwordResetReq)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

func TestChangePasswordAPI(t *testing.T) {
	changePasswordReq := &ChangePasswordRequest{
		CurrentPassword:    constant.DummyPassword,
		NewPassword:        constant.DummyPassword,
		ConfirmNewPassword: constant.DummyPassword,
	}

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

		err = repo.ChangePasswordAPI(constant.DummyMemberId, changePasswordReq)

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

		err := repo.ChangePasswordAPI(constant.DummyMemberId, changePasswordReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), constant.ErrFailedMarshalReq)
	})

	t.Run("NewRequestError", func(t *testing.T) {
		mockClient := new(MockClient)
		repo := NewRepository(&config.Config{
			Env: &config.Environment{AifcoreHost: constant.MockInvalidHost},
		}, mockClient, nil)

		err := repo.ChangePasswordAPI(constant.DummyMemberId, changePasswordReq)
		assert.Error(t, err)
	})

	t.Run("HTTPRequestError", func(t *testing.T) {
		expectedErr := errors.New(constant.ErrHTTPReqFailed)

		repo, mockClient := setupMockRepo(t, nil, expectedErr)

		err := repo.ChangePasswordAPI(constant.DummyMemberId, changePasswordReq)

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

		err := repo.ChangePasswordAPI(constant.DummyMemberId, changePasswordReq)
		assert.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}

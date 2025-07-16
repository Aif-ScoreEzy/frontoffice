package transaction

import (
	"front-office/app/config"
	"front-office/common/constant"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
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

package transaction

// import (
// 	"front-office/app/config"
// 	"front-office/common/constant"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func initMockConfig(mockServerURL string) *config.Config {
// 	return &config.Config{
// 		Env: &config.Environment{
// 			AifcoreHost: mockServerURL,
// 		},
// 	}
// }

// type MockClient struct {
// 	mock.Mock
// }

// func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
// 	args := m.Called(req)
// 	return args.Get(0).(*http.Response), args.Error(1)
// }

// const dummyResponseBody = `{
// 	"success": true,
// 	"data": [
// 		{
// 			"log_trx_id": 1,
// 			"trx_id": "1",
// 			"user_id": 1,
// 			"company_id": 1,
// 			"ip_client": "",
// 			"product_id": 0,
// 			"status": "",
// 			"success": false,
// 			"message": "",
// 			"probability_to_default": "",
// 			"grade": "",
// 			"loan_no": "",
// 			"data": null,
// 			"duration": 0,
// 			"created_at": "2024-09-29T19:53:34.819141+07:00"
// 		}
// 	],
// 	"message": "",
// 	"meta": {
// 		"total": 1,
// 		"page": 1,
// 		"total_pages": 1,
// 		"visible": 1,
// 		"start_data": 1,
// 		"end_data": 1,
// 		"size": 10,
// 		"message": "Success"
// 	}
// }`
// const dummyDate = "2024-09-25"

// func TestNewRepository(t *testing.T) {
// 	// Create mock instances of gorm.DB and config.Config
// 	mockConfig := &config.Config{}

// 	// Act - Call NewRepository
// 	mockClient := new(MockClient)
// 	repo := NewRepository(mockConfig, mockClient)

// 	// Assert - Ensure that repo is not nil
// 	assert.NotNil(t, repo)

// 	// Assert - Ensure that repo is of the correct type (*repository)
// 	_, ok := repo.(*repository)
// 	assert.True(t, ok, "Expected *repository, got something else")

// 	// Assert - Ensure that the DB and Cfg fields are correctly assigned
// 	assert.Equal(t, mockConfig, repo.(*repository).cfg)
// }

// func TestFindAllTransactionLogs(t *testing.T) {
// 	// Mock server to simulate the external API
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Validate the request URL and method
// 		expectedPath := "/api/core/logging/transaction/scoreezy/list"
// 		assert.Equal(t, expectedPath, r.URL.Path)
// 		assert.Equal(t, http.MethodGet, r.Method)

// 		// Mock response body
// 		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(dummyResponseBody))
// 	}))
// 	defer mockServer.Close()

// 	// Initialize mock config with the mock server's URL
// 	mockConfig := initMockConfig(mockServer.URL)
// 	mockClient := new(MockClient)

// 	repo := &repository{
// 		cfg:    mockConfig,
// 		client: mockClient,
// 	}

// 	resp, err := repo.CallScoreezyLogsAPI()

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NoError(t, err)

// }

// func TestFindAllTransactionLogsByDate(t *testing.T) {
// 	// Mock server to simulate the external API
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Validate the request URL and method
// 		expectedPath := "/api/core/logging/transaction/scoreezy/by"
// 		assert.Equal(t, expectedPath, r.URL.Path)
// 		assert.Equal(t, http.MethodGet, r.Method)

// 		// Validate query parameters
// 		query := r.URL.Query()
// 		assert.Equal(t, "1", query.Get("company_id"))
// 		assert.Equal(t, dummyDate, query.Get("date"))

// 		// Mock response body
// 		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(dummyResponseBody))
// 	}))
// 	defer mockServer.Close()

// 	// Initialize mock config with the mock server's URL
// 	mockConfig := initMockConfig(mockServer.URL)
// 	mockClient := new(MockClient)

// 	repo := &repository{
// 		cfg:    mockConfig,
// 		client: mockClient,
// 	}

// 	// Act - Call the method with companyId and date
// 	companyId := "1"
// 	date := dummyDate
// 	resp, err := repo.CallScoreezyLogsByDateAPI(companyId, date)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NoError(t, err)
// }

// func TestFindAllTransactionLogsByRangeDate(t *testing.T) {
// 	// Mock server to simulate the external API
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Validate the request URL and method
// 		expectedPath := "/api/core/logging/transaction/scoreezy/range"
// 		assert.Equal(t, expectedPath, r.URL.Path)
// 		assert.Equal(t, http.MethodGet, r.Method)

// 		// Validate query parameters
// 		query := r.URL.Query()
// 		assert.Equal(t, "1", query.Get("company_id"))
// 		assert.Equal(t, dummyDate, query.Get("date_start"))
// 		assert.Equal(t, dummyDate, query.Get("date_end"))

// 		// Mock response body
// 		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(dummyResponseBody))
// 	}))
// 	defer mockServer.Close()

// 	// Initialize mock config with the mock server's URL
// 	mockConfig := initMockConfig(mockServer.URL)
// 	mockClient := new(MockClient)

// 	repo := &repository{
// 		cfg:    mockConfig,
// 		client: mockClient,
// 	}

// 	// Act - Call the method with companyId and date
// 	companyId := "1"
// 	startDate := dummyDate
// 	endDate := dummyDate
// 	resp, err := repo.CallScoreezyLogsByRangeDateAPI(companyId, startDate, endDate)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NoError(t, err)
// }

// func TestFindAllTransactionLogsByMonth(t *testing.T) {
// 	// Mock server to simulate the external API
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		expectedPath := "/api/core/logging/transaction/scoreezy/month"
// 		assert.Equal(t, expectedPath, r.URL.Path)
// 		assert.Equal(t, http.MethodGet, r.Method)

// 		query := r.URL.Query()
// 		assert.Equal(t, "1", query.Get("company_id"))
// 		assert.Equal(t, dummyDate, query.Get("month"))

// 		// Mock response body
// 		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(dummyResponseBody))
// 	}))
// 	defer mockServer.Close()

// 	mockConfig := initMockConfig(mockServer.URL)
// 	mockClient := new(MockClient)

// 	repo := &repository{
// 		cfg:    mockConfig,
// 		client: mockClient,
// 	}

// 	companyId := "1"
// 	date := dummyDate
// 	resp, err := repo.CallScoreezyLogsByMonthAPI(companyId, date)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.NoError(t, err)
// }

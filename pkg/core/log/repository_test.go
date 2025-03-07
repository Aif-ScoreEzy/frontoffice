package log

import (
	"front-office/app/config"
	"front-office/common/constant"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initMockConfig(mockServerURL string) *config.Config {
	return &config.Config{
		Env: &config.Environment{
			AifcoreHost: mockServerURL,
		},
	}
}

const dummyResponseBody = `{
	"success": true,
	"data": [
		{
			"log_trx_id": 1,
			"trx_id": "1",
			"user_id": 1,
			"company_id": 1,
			"ip_client": "",
			"product_id": 0,
			"status": "",
			"success": false,
			"message": "",
			"probability_to_default": "",
			"grade": "",
			"loan_no": "",
			"data": null,
			"duration": 0,
			"created_at": "2024-09-29T19:53:34.819141+07:00"
		}
	],
	"message": "",
	"meta": {
		"total": 1,
		"page": 1,
		"total_pages": 1,
		"visible": 1,
		"start_data": 1,
		"end_data": 1,
		"size": 10,
		"message": "Success"
	}
}`
const dummyDate = "2024-09-25"

func TestNewRepository(t *testing.T) {
	// Create mock instances of gorm.DB and config.Config
	mockDB := &gorm.DB{}
	mockConfig := &config.Config{}

	// Act - Call NewRepository
	repo := NewRepository(mockDB, mockConfig)

	// Assert - Ensure that repo is not nil
	assert.NotNil(t, repo)

	// Assert - Ensure that repo is of the correct type (*repository)
	_, ok := repo.(*repository)
	assert.True(t, ok, "Expected *repository, got something else")

	// Assert - Ensure that the DB and Cfg fields are correctly assigned
	assert.Equal(t, mockDB, repo.(*repository).DB)
	assert.Equal(t, mockConfig, repo.(*repository).Cfg)
}

func TestFindAllTransactionLogs(t *testing.T) {
	// Mock server to simulate the external API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request URL and method
		expectedPath := "/api/core/logging/transaction/list"
		assert.Equal(t, expectedPath, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		// Mock response body
		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(dummyResponseBody))
	}))
	defer mockServer.Close()

	// Initialize mock config with the mock server's URL
	mockConfig := initMockConfig(mockServer.URL)

	repo := &repository{
		DB:  &gorm.DB{},
		Cfg: mockConfig,
	}

	resp, err := repo.FetchLogTransactions()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expectedBody := dummyResponseBody
	assert.JSONEq(t, expectedBody, string(body))
}

func TestFindAllTransactionLogsByDate(t *testing.T) {
	// Mock server to simulate the external API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request URL and method
		expectedPath := "/api/core/logging/transaction/by"
		assert.Equal(t, expectedPath, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		// Validate query parameters
		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("company_id"))
		assert.Equal(t, dummyDate, query.Get("date"))

		// Mock response body
		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(dummyResponseBody))
	}))
	defer mockServer.Close()

	// Initialize mock config with the mock server's URL
	mockConfig := initMockConfig(mockServer.URL)

	// Create repository with the mock config
	repo := &repository{
		DB:  &gorm.DB{},
		Cfg: mockConfig,
	}

	// Act - Call the method with companyId and date
	companyId := "1"
	date := dummyDate
	resp, err := repo.FetchLogTransactionsByDate(companyId, date)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Optionally, check the body content
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expectedBody := dummyResponseBody

	assert.JSONEq(t, expectedBody, string(body))
}

func TestFindAllTransactionLogsByRangeDate(t *testing.T) {
	// Mock server to simulate the external API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the request URL and method
		expectedPath := "/api/core/logging/transaction/range"
		assert.Equal(t, expectedPath, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		// Validate query parameters
		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("company_id"))
		assert.Equal(t, dummyDate, query.Get("date_start"))
		assert.Equal(t, dummyDate, query.Get("date_end"))

		// Mock response body
		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(dummyResponseBody))
	}))
	defer mockServer.Close()

	// Initialize mock config with the mock server's URL
	mockConfig := initMockConfig(mockServer.URL)

	// Create repository with the mock config
	repo := &repository{
		DB:  &gorm.DB{},
		Cfg: mockConfig,
	}

	// Act - Call the method with companyId and date
	companyId := "1"
	startDate := dummyDate
	endDate := dummyDate
	resp, err := repo.FetchLogTransactionsByRangeDate(companyId, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Optionally, check the body content
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expectedBody := dummyResponseBody

	assert.JSONEq(t, expectedBody, string(body))
}

func TestFindAllTransactionLogsByMonth(t *testing.T) {
	// Mock server to simulate the external API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/core/logging/transaction/month"
		assert.Equal(t, expectedPath, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		query := r.URL.Query()
		assert.Equal(t, "1", query.Get("company_id"))
		assert.Equal(t, dummyDate, query.Get("month"))

		// Mock response body
		w.Header().Set(constant.HeaderContentType, constant.HeaderApplicationJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(dummyResponseBody))
	}))
	defer mockServer.Close()

	mockConfig := initMockConfig(mockServer.URL)

	repo := &repository{
		DB:  &gorm.DB{},
		Cfg: mockConfig,
	}

	companyId := "1"
	date := dummyDate
	resp, err := repo.FetchLogTransactionsByMonth(companyId, date)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	expectedBody := dummyResponseBody
	assert.JSONEq(t, expectedBody, string(body))
}

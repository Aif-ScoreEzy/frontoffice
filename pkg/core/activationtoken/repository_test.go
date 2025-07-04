package activationtoken

// import (
// 	"fmt"
// 	"front-office/app/config"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// // Define the struct for expected response data
// type ActivationTokenResponse struct {
// 	Success bool `json:"success"`
// 	Data    struct {
// 		ID        string    `json:"id"`
// 		Token     string    `json:"token"`
// 		MemberID  int       `json:"member_id"`
// 		CreatedAt time.Time `json:"created_at"`
// 	} `json:"data"`
// 	Message string `json:"message"`
// 	Meta    struct {
// 		Message string `json:"message"`
// 	} `json:"meta"`
// }

// func initMockConfig(mockServerURL string) *config.Config {
// 	return &config.Config{
// 		Env: &config.Environment{
// 			AifcoreHost: mockServerURL,
// 		},
// 	}
// }

// func TestFindOneActivationTokenBytoken(t *testing.T) {
// 	// Mock server to simulate the external API
// 	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiNTQyYzRmNjItZTk5MC00Y2VjLWFlNDAtZjZlZTU5YzlhNGU4IiwiZXhwIjoxNzE2ODg3ODA4LCJ0aWVyX2xldmVsIjoxLCJ1c2VyX2lkIjoiYjY4NTQ0MWUtNWFkNS00NTRkLWEzNWMtMmI2MGZhZTY2Y2M0In0.5x9R9tYQ7dtYdR5iUUlfPRA0YJKJVN1260gMBFcdH1w"
// 		expectedPath := fmt.Sprintf("/api/core/member/activation-tokens/%s", expectedToken)

// 		assert.Equal(t, expectedPath, r.URL.Path) // Check if the path is correct
// 		assert.Equal(t, http.MethodGet, r.Method) // Check if the method is GET

// 		// Mock response body
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write([]byte(`{
// 			"success": true,
// 			"data": {
// 				"id": "d5f6975e-4537-4025-a677-7abcb98a5aec",
// 				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiNTQyYzRmNjItZTk5MC00Y2VjLWFlNDAtZjZlZTU5YzlhNGU4IiwiZXhwIjoxNzE2ODg3ODA4LCJ0aWVyX2xldmVsIjoxLCJ1c2VyX2lkIjoiYjY4NTQ0MWUtNWFkNS00NTRkLWEzNWMtMmI2MGZhZTY2Y2M0In0.5x9R9tYQ7dtYdR5iUUlfPRA0YJKJVN1260gMBFcdH1w",
// 				"member_id": 5,
// 				"created_at": "2024-09-18T17:05:22.348406+07:00"
// 			},
// 			"message": "",
// 			"meta": {
// 				"message": "Success to get activation token"
// 			}
// 		}`))
// 	}))
// 	defer mockServer.Close()

// 	// Initialize mock config with the mock server's URL
// 	mockConfig := initMockConfig(mockServer.URL)

// 	// Create repository with the mock config
// 	repo := &repository{
// 		cfg: mockConfig,
// 	}

// 	// Act - Call the method with the provided token
// 	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiNTQyYzRmNjItZTk5MC00Y2VjLWFlNDAtZjZlZTU5YzlhNGU4IiwiZXhwIjoxNzE2ODg3ODA4LCJ0aWVyX2xldmVsIjoxLCJ1c2VyX2lkIjoiYjY4NTQ0MWUtNWFkNS00NTRkLWEzNWMtMmI2MGZhZTY2Y2M0In0.5x9R9tYQ7dtYdR5iUUlfPRA0YJKJVN1260gMBFcdH1w"
// 	resp, err := repo.CallGetActivationTokenAPI(token)

// 	// Assert
// 	assert.NoError(t, err) // Check that no error occurred
// 	assert.NotNil(t, resp) // Ensure that the response is not nil
// 	// assert.Equal(t, http.StatusOK, resp.StatusCode) // Check that the status code is 200 OK

// 	// Optionally, check the body content (here we just validate the response, but you could unmarshal and check JSON)
// 	// defer resp.Body.Close()
// 	// body, err := io.ReadAll(resp.Body)
// 	assert.NoError(t, err)

// 	var actualResponse ActivationTokenResponse
// 	// err = json.Unmarshal(body, &actualResponse)
// 	assert.NoError(t, err)

// 	// Compare the actual response to the expected response
// 	expectedResponse := ActivationTokenResponse{
// 		Success: true,
// 		Data: struct {
// 			ID        string    `json:"id"`
// 			Token     string    `json:"token"`
// 			MemberID  int       `json:"member_id"`
// 			CreatedAt time.Time `json:"created_at"`
// 		}{
// 			ID:        "d5f6975e-4537-4025-a677-7abcb98a5aec",
// 			Token:     token,
// 			MemberID:  5,
// 			CreatedAt: time.Date(2024, 9, 18, 17, 5, 22, 348406000, time.Local),
// 		},
// 		Meta: struct {
// 			Message string `json:"message"`
// 		}{
// 			Message: "Success to get activation token",
// 		},
// 	}

// 	// Validate the response content
// 	assert.Equal(t, expectedResponse.Success, actualResponse.Success)
// 	assert.Equal(t, expectedResponse.Data.ID, actualResponse.Data.ID)
// 	assert.Equal(t, expectedResponse.Data.Token, actualResponse.Data.Token)
// 	assert.Equal(t, expectedResponse.Data.MemberID, actualResponse.Data.MemberID)
// 	assert.Equal(t, expectedResponse.Data.CreatedAt, actualResponse.Data.CreatedAt)
// 	assert.Equal(t, expectedResponse.Meta.Message, actualResponse.Meta.Message)
// }

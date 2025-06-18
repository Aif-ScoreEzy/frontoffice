package taxcompliancestatus

import (
	"bytes"
	"front-office/app/config"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCallTaxComplianceStatusAPI(t *testing.T) {
	cfg := &config.Config{
		Env: &config.Environment{
			ProductCatalogHost: "http://mock-host",
		},
	}
	mockClient := new(MockClient)
	repo := NewRepository(cfg, mockClient)

	expectedBody := `{
		"message": "Succeed to Request Data.",
		"success": true,
		"input": {"npwp": "092542823407000"},
		"data": {
			"nama": "AIFORESEE INOVASI SKOR",
  			"alamat": "AIA CENTRAL BUILDING LANTAI 28, JL. JENDERAL SUDIRMAN KAV. 48 A   005 004",
   			"status": "Reported"
		},
		"pricing_strategy": "PAY",
		"transaction_id": "9c6b46c9-e3be-4c90-a5e3-894b26432e0b",
		"datetime": "2025-06-01 05:17:53"
	}`

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(expectedBody))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	mockClient.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(mockResp, nil)

	req := &taxComplianceStatusRequest{
		Npwp: "092542823407000",
	}

	resp, err := repo.CallTaxComplianceStatusAPI("test-api-key", "1", req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

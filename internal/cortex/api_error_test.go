package cortex_test

import (
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiErrorString(t *testing.T) {
	testApiError := cortex.ApiError{
		Details:           "test-details",
		GatewayHttpStatus: 200,
		HttpStatus:        200,
		Message:           "test-message",
		RequestId:         "test-request-id",
		Type:              "test-type",
	}
	assert.Contains(t, testApiError.String(), "error: test-type")
	assert.Contains(t, testApiError.String(), "detail: test-details")
	assert.Contains(t, testApiError.String(), "message: test-message")
	assert.Contains(t, testApiError.String(), "status: 200")
}

package cortex_test

import (
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiErrorString(t *testing.T) {
	testApiError := cortex.ApiError{
		Details:           "test-details",
		GatewayHttpStatus: 201,
		HttpStatus:        200,
		Message:           "test-message",
		RequestId:         "test-request-id",
		Type:              "test-type",
	}
	assert.Contains(t, testApiError.String(), "type: test-type")
	assert.Contains(t, testApiError.String(), "details: test-details")
	assert.Contains(t, testApiError.String(), "message: test-message")
	assert.Contains(t, testApiError.String(), "http status: 200")
	assert.Contains(t, testApiError.String(), "gateway status: 201")
	assert.Contains(t, testApiError.String(), "requestId: test-request-id")
}

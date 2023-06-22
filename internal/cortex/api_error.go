package cortex

import (
	"errors"
	"fmt"
)

var (
	ApiErrorNotFound     = errors.New("not found")
	ApiErrorUnauthorized = errors.New("unauthorized")
)

type ApiError struct {
	Details           string `json:"details"`
	GatewayHttpStatus int    `json:"gatewayHttpStatus"`
	HttpStatus        int    `json:"httpStatus"`
	Message           string `json:"message"`
	RequestId         string `json:"requestId"`
	Type              string `json:"type"`
}

func (err ApiError) String() string {
	str := ""
	if err.Type != "" {
		str += fmt.Sprintf("type: %s\n", err.Type)
	}

	if err.Details != "" {
		str += fmt.Sprintf("details: %s\n", err.Details)
	}

	if err.Message != "" {
		str += fmt.Sprintf("message: %s\n", err.Message)
	}

	if err.HttpStatus > 0 {
		str += fmt.Sprintf("http status: %d\n", err.HttpStatus)
	}

	if err.GatewayHttpStatus > 0 {
		str += fmt.Sprintf("gateway status: %d\n", err.GatewayHttpStatus)
	}

	str += fmt.Sprintf("requestId: %s\n", err.RequestId)

	return str
}

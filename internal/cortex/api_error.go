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
		str += fmt.Sprintf("error: %s\n", err.Type)
	}

	if err.Details != "" {
		str += fmt.Sprintf("detail: %s\n", err.Details)
	}

	if err.Message != "" {
		str += fmt.Sprintf("message: %s", err.Message)
	}

	if err.HttpStatus > 0 {
		str += fmt.Sprintf("status: %d", err.HttpStatus)
	}

	return str
}

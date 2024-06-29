package common

import (
	"net/http"

	"github.com/pkxro/squid/internal/model"
)

var (
	// SquidBuildtimeError is an error used when there are buildtime issues
	SquidBuildtimeError = StandardSentinelError{
		Status:  http.StatusServiceUnavailable,
		Message: "unknown error occurred when starting the server",
	}

	// SquidBadRequestError is an error used to show a request was invalid
	SquidBadRequestError = StandardSentinelError{
		Status:  http.StatusBadRequest,
		Message: "bad request",
	}

	// SquidIdempotencyError is an error used to show a request was invalid
	SquidIdempotencyError = StandardSentinelError{
		Status:  http.StatusForbidden,
		Message: "idempotency key is invalid",
	}

	// SquidNotFoundError is an error used to show a request was invalid
	SquidNotFoundError = StandardSentinelError{
		Status:  http.StatusNotFound,
		Message: "not found",
	}

	// SquidInternalError is an error used to show a request was invalid
	SquidInternalError = StandardSentinelError{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}
)

// APIError is an interface for accessing and implementing a client error
type APIError interface {
	APIError() (int, string)
}

// StandardSentinelError is a struct for returning embedded runtime errors
type StandardSentinelError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// WrapError is used to create a public facing wrapped errors that
// includes the `sentinelError` which is an internal or service application errors
func WrapError(errMsg string, sentinel StandardSentinelError, av model.APIVersion) model.StandardErrorResponse {
	return model.StandardErrorResponse{InternalErrMsg: errMsg, Error: sentinel, APIVersion: av}
}

// WrapAPIError is used to create a public facing wrapped errors that
// includes the `sentinelError` which is an internal or service application errors
func WrapAPIError(errMsg string, sentinel StandardSentinelError, av model.APIVersion) (int, model.StandardErrorResponse) {
	return sentinel.Status, model.StandardErrorResponse{InternalErrMsg: errMsg, Error: sentinel, APIVersion: av}
}

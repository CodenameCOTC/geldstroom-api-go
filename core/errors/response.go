package errorsresponse

import "net/http"

type ErrorResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

type ValidationErrorResponse struct {
	ErrorCode string            `json:"errorCode"`
	Message   string            `json:"message"`
	Error     map[string]string `json:"error"`
}

type BadRequestResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func InternalServerError(message string) ErrorResponse {
	if message == "" {
		message = "We encountered an error while processing your request."
	}

	return ErrorResponse{
		ErrorCode: http.StatusInternalServerError,
		Message:   message,
	}
}

func ValidationError(errorCode string, message error, error map[string]string) ValidationErrorResponse {
	return ValidationErrorResponse{
		ErrorCode: errorCode,
		Message:   message.Error(),
		Error:     error,
	}
}

func NotFound(message string) ErrorResponse {
	if message == "" {
		message = "The requested resource was not found."
	}

	return ErrorResponse{
		ErrorCode: http.StatusNotFound,
		Message:   message,
	}
}

func Unauthorized(message string) ErrorResponse {
	if message == "" {
		message = "You are not authenticated to perform the requested action."
	}

	return ErrorResponse{
		ErrorCode: http.StatusUnauthorized,
		Message:   message,
	}
}

func Forbidden(message string) ErrorResponse {
	if message == "" {
		message = "You are not authorized to perform the requested action."
	}

	return ErrorResponse{
		ErrorCode: http.StatusForbidden,
		Message:   message,
	}
}

func InvalidQuery(errorCode string, message error, error map[string]string) ValidationErrorResponse {
	return ValidationErrorResponse{
		ErrorCode: errorCode,
		Message:   message.Error(),
		Error:     error,
	}
}

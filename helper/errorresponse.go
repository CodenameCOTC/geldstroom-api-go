package helper

type ErrorResponse struct {
	Message   string            `json:"message"`
	ErrorCode string            `json:"errorCode"`
	Error     map[string]string `json:"error"`
}

var InternalServerError = map[string]string{
	"message": "Internal Server Error",
}

var Unauthorized = map[string]string{
	"message": "Unauthorized",
}

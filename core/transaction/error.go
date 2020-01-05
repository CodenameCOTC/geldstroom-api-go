package transaction

import "errors"

var (
	ErrValidationFailed        = errors.New("Validation Failed")
	ErrValidationFailedCode    = "TRANSACTION_0001"
	ErrInvalidQueryCode        = "TRANSACTION_0002"
	ErrTransactionNotFound     = errors.New("Transaction Not Found")
	ErrTransactionNotFoundCode = "TRANSACTION_0004"
)

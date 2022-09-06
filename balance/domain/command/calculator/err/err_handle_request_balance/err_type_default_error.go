package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorTypeDefaultError     = error_base.Error
	ErrorTypeDefaultErrorType = error_base.ErrorBase
)

const (
	ERROR_CODE_DEFAULT_ERROR      = 500
	ERROR_MESSAGE_DEFAULT_ERROR   = "Have error, please try again"
	ERROR_SIGNATURE_DEFAULT_ERROR = "err_handle_request_balance_ERROR_SIGNATURE_DEFAULT_ERROR"
)

func NewErrorDefaultError() ErrorTypeDefaultError {
	return &ErrorTypeDefaultErrorType{
		Code:      ERROR_CODE_DEFAULT_ERROR,
		Message:   ERROR_MESSAGE_DEFAULT_ERROR,
		Signature: ERROR_SIGNATURE_DEFAULT_ERROR,
	}
}

func IsErrorDefaultError(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_DEFAULT_ERROR)
}

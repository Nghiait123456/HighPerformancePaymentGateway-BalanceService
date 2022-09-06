package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorTypeRequestBalanceNotValid     = error_base.Error
	ErrorTypeRequestBalanceNotValidType = error_base.ErrorBase
)

const (
	ERROR_CODE_TYPE_REQUEST_BALANCE_NOT_VALID      = 5
	ERROR_MESSAGE_TYPE_REQUEST_BALANCE_NOT_VALID   = "type request balance is not valid"
	ERROR_SIGNATURE_TYPE_REQUEST_BALANCE_NOT_VALID = "err_handle_request_balance_ERROR_SIGNATURE_TYPE_REQUEST_BALANCE_NOT_VALID"
)

func NewErrorTypeRequestBalanceNotValid() ErrorTypeRequestBalanceNotValid {
	return &ErrorTypeRequestBalanceNotValidType{
		Code:      ERROR_CODE_TYPE_REQUEST_BALANCE_NOT_VALID,
		Message:   ERROR_MESSAGE_TYPE_REQUEST_BALANCE_NOT_VALID,
		Signature: ERROR_SIGNATURE_TYPE_REQUEST_BALANCE_NOT_VALID,
	}
}

func IsErrorTypeRequestBalanceNotValid(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_TYPE_REQUEST_BALANCE_NOT_VALID)
}

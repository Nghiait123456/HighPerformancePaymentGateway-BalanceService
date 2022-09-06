package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorBalanceIsFull     = error_base.Error
	ErrorBalanceIsFullType = error_base.ErrorBase
)

const (
	ERROR_CODE_BALANCE_IS_FULL      = 2
	ERROR_MESSAGE_BALANCE_IS_FULL   = "Balance is full, please try again late"
	ERROR_SIGNATURE_BALANCE_IS_FULL = "err_handle_request_balance_ERROR_SIGNATURE_BALANCE_IS_FULL"
)

func NewErrorBalanceIsFull() ErrorBalanceIsFull {
	return &ErrorBalanceIsFullType{
		Code:      ERROR_CODE_BALANCE_IS_FULL,
		Message:   ERROR_MESSAGE_BALANCE_IS_FULL,
		Signature: ERROR_SIGNATURE_BALANCE_IS_FULL,
	}
}

func IsErrorBalanceIsFull(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_BALANCE_IS_FULL)
}

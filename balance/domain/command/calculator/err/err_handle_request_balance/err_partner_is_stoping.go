package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorPartnerIsStoppingForEmergenceStop     = error_base.Error
	ErrorPartnerIsStoppingForEmergenceStopType = error_base.ErrorBase
)

const (
	ERROR_CODE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP      = 1
	ERROR_MESSAGE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP   = "Partner is stopping for emergence stop, please try again"
	ERROR_SIGNATURE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP = "err_handle_request_balance_ERROR_SIGNATURE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP"
)

func NewErrorPartnerStoppingForEmergenceStop() ErrorPartnerIsStoppingForEmergenceStop {
	return &ErrorPartnerIsStoppingForEmergenceStopType{
		Code:      ERROR_CODE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP,
		Message:   ERROR_MESSAGE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP,
		Signature: ERROR_SIGNATURE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP,
	}
}

func IsErrorPartnerStoppingForEmergenceStop(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_PARTNER_IS_STOPPING_FOR_EMERGENCE_STOP)
}

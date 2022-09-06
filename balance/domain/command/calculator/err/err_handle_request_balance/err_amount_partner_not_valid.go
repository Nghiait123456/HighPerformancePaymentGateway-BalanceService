package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorAmountPartnerNotValid     = error_base.Error
	ErrorAmountPartnerNotValidType = error_base.ErrorBase
)

const (
	ERROR_CODE_AMOUNT_PARTNER_NOT_VALID      = 3
	ERROR_MESSAGE_AMOUNT_PARTNER_NOT_VALID   = "Amount partner is not valid"
	ERROR_SIGNATURE_AMOUNT_PARTNER_NOT_VALID = "err_handle_request_balance_ERROR_SIGNATURE_AMOUNT_PARTNER_NOT_VALID"
)

func NewErrorAmountPartnerIsNotValid() ErrorAmountPartnerNotValid {
	return &ErrorAmountPartnerNotValidType{
		Code:      ERROR_CODE_AMOUNT_PARTNER_NOT_VALID,
		Message:   ERROR_MESSAGE_AMOUNT_PARTNER_NOT_VALID,
		Signature: ERROR_SIGNATURE_AMOUNT_PARTNER_NOT_VALID,
	}
}

func IsErrorAmountPartnerNotValid(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_AMOUNT_PARTNER_NOT_VALID)
}

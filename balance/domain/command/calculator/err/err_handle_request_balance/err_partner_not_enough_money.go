package err_handle_request_balance

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorPartnerNotEnoughMoney     = error_base.Error
	ErrorPartnerNotEnoughMoneyType = error_base.ErrorBase
)

const (
	ERROR_CODE_PARTNER_NOT_ENOUGH_MONEY      = 4
	ERROR_MESSAGE_PARTNER_NOT_ENOUGH_MONEY   = "Partner is not enough money, please try again"
	ERROR_SIGNATURE_PARTNER_NOT_ENOUGH_MONEY = "err_handle_request_balance_ERROR_SIGNATURE_PARTNER_NOT_ENOUGH_MONEY"
)

func NewErrorPartnerNotEnoughMoney() ErrorPartnerNotEnoughMoney {
	return &ErrorPartnerNotEnoughMoneyType{
		Code:      ERROR_CODE_PARTNER_NOT_ENOUGH_MONEY,
		Message:   ERROR_MESSAGE_PARTNER_NOT_ENOUGH_MONEY,
		Signature: ERROR_SIGNATURE_PARTNER_NOT_ENOUGH_MONEY,
	}
}

func IsErrorPartnerNotEnoughMoney(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_SIGNATURE_PARTNER_NOT_ENOUGH_MONEY)
}

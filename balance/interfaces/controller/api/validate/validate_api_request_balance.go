package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_request"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/validate"
)

type (
	ValidateApiRequestBalance struct {
		BV  validate.ValidateBaseInterface
		dto dto_api_request.RequestBalanceDto
	}
)

const (
	MIN_AMOUNT = 10000
	MAX_AMOUNT = 5000000
)

func (v *ValidateApiRequestBalance) Init(fl validator.FieldLevel) {
	v.BV.ResignValidateCustom("minAmount", v.minAmount)
	v.BV.ResignValidateCustom("maxAmount", v.maxAmount)
	v.BV.ResignValidateCustom("partnerExist", v.partnerExist)
	v.BV.ResignValidateCustom("partnerActive", v.partnerActive)
	v.BV.ResignValidateCustom("orderExist", v.orderExist)
	v.BV.ResignValidateCustom("orderStatusValid", v.orderStatusValid)

	message := make(validate.MapMessage)
	message["minAmount"] = "Amount is less than min allow"
	message["maxAmount"] = "Amount is  greater max allow"
	message["partnerExist"] = "Partner is not exist"
	message["partnerActive"] = "Partner is not Active"
	message["orderExist"] = "Order is not exist"
	message["orderStatusValid"] = "Order status is not valid"
	v.BV.SetMessageForRule(message)
}

func (v ValidateApiRequestBalance) minAmount(fl validator.FieldLevel) bool {
	return fl.Field().Uint() > MIN_AMOUNT
}

func (v ValidateApiRequestBalance) maxAmount(fl validator.FieldLevel) bool {
	return fl.Field().Uint() < MAX_AMOUNT
}

func (v ValidateApiRequestBalance) partnerExist(fl validator.FieldLevel) bool {
	//todo check
	return true
}

func (v ValidateApiRequestBalance) partnerActive(fl validator.FieldLevel) bool {
	//todo check
	return true
}

func (v ValidateApiRequestBalance) orderExist(fl validator.FieldLevel) bool {
	//todo check
	return true
}

func (v ValidateApiRequestBalance) orderStatusValid(fl validator.FieldLevel) bool {
	//todo check
	return true
}

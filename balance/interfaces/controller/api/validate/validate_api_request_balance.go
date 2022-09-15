package validate

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_request"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_response"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/validate"
	"github.com/high-performance-payment-gateway/balance-service/balance/value_object"
	log "github.com/sirupsen/logrus"
)

type (
	ValidateApiRequestBalance struct {
		VB  validate.ValidateBaseInterface
		Dto dto_api_request.RequestBalanceDto
	}
)

const (
	MIN_AMOUNT = 10000
	MAX_AMOUNT = 5000000
)

func (v *ValidateApiRequestBalance) Init() {
	v.VB.ResignValidateCustom("minAmount", v.minAmount)
	v.VB.ResignValidateCustom("maxAmount", v.maxAmount)
	v.VB.ResignValidateCustom("partnerExist", v.partnerExist)
	v.VB.ResignValidateCustom("partnerActive", v.partnerActive)
	v.VB.ResignValidateCustom("orderExist", v.orderExist)
	v.VB.ResignValidateCustom("orderStatusValid", v.orderStatusValid)
	v.VB.ResignValidateCustom("typeRequestExist", v.typeRequestExist)

	message := make(validate.MapMessage)
	message["minAmount"] = "Amount is less than min allow"
	message["maxAmount"] = "Amount is  greater max allow"
	message["partnerExist"] = "Partner is not exist"
	message["partnerActive"] = "Partner is not Active"
	message["orderExist"] = "Order is not exist"
	message["orderStatusValid"] = "Order status is not valid"
	message["typeRequestExist"] = v.MessageErrorTypeRequestNotValid()
	v.VB.SetMessageForRule(message)
}

// return struct response, error
func (v *ValidateApiRequestBalance) Validate() (dto_api_response.ResponseRequestBalanceDto, error) {
	errV := v.VB.Validate().Struct(v.Dto)
	if errV != nil {
		message, errCE := v.VB.ConvertErrorValidate(errV)
		if errCE != nil {
			fmt.Println("invalidate error")
			res := dto_api_response.ResponseRequestBalanceDto{
				HttpCode:    http_status.StatusBadRequest,
				Status:      dto_api_response.STATUS_ERROR,
				Code:        http_status.StatusBadRequest,
				Message:     "Param is invalid format, please check and try again",
				ErrorDetail: errV.Error(),
			}
			return res, errV
		}

		fmt.Println("message =", message)

		// show message
		errSE, detail := v.VB.ShowErrors(message, v.CustomShowError)
		fmt.Println("detail =", detail)
		if errSE != nil {
			messageErr := fmt.Sprintf("ShowErrors validate error: %s", errSE.Error())
			log.WithFields(log.Fields{
				"errMessage": errSE.Error(),
			}).Error("")
			panic(messageErr)
		}

		res := dto_api_response.ResponseRequestBalanceDto{
			HttpCode:    http_status.StatusBadRequest,
			Status:      dto_api_response.STATUS_ERROR,
			Code:        http_status.StatusBadRequest,
			Message:     "Param missing or invalid format, please check and try again",
			ErrorDetail: detail,
		}
		fmt.Println("res", res)

		return res, errors.New("Validate has errors")
	}

	return dto_api_response.ResponseRequestBalanceDto{}, nil
}

func (v ValidateApiRequestBalance) CustomShowError(mE validate.MessageErrors) (error, interface{}) {
	ListES := validate.ListErrorsDefaultShow{}

	for _, v := range mE {
		oneErr := validate.OneErrorDefaultShow{
			Field:      v.Field,
			Rule:       v.Rule,
			Message:    v.Message,
			ParamRule:  v.ParamRule,
			ValueError: v.ValueError,
		}

		ListES[v.Field] = oneErr
	}

	showE := validate.DefaultShowError{
		ListError: ListES,
	}

	return nil, showE
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

func (v ValidateApiRequestBalance) typeRequestExist(fl validator.FieldLevel) bool {
	t := value_object.NewTypeRequestBalance()
	return t.TypeRequestExist(fl.Field().String())
}

func (v ValidateApiRequestBalance) MessageErrorTypeRequestNotValid() string {
	t := value_object.NewTypeRequestBalance()
	allType := fmt.Sprintf("%v", t.AllType)
	message := fmt.Sprintf("typeRequest is not valid, in list: %s", allType)

	return message
}

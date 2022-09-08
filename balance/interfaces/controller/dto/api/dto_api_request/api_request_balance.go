package dto_api_request

import (
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_response"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
	log "github.com/sirupsen/logrus"
)

type (
	RequestBalanceDtoRaw struct {
		AmountRequest         any `json:"AmountRequest" ml:"AmountRequest" form:"AmountRequest"`
		PartnerCode           any `json:"PartnerCode" ml:"PartnerCode" form:"PartnerCode"`
		PartnerIdentification any `json:"PartnerIdentification" ml:"PartnerIdentification" form:"PartnerIdentification"`
		OrderID               any `json:"OrderID" ml:"OrderID" form:"OrderID"`
		// create order, update amount when partner recharge
		TypeRequest any `json:"TypeRequest" ml:"TypeRequest" form:"TypeRequest"`
	}

	RequestBalanceDto struct {
		AmountRequest         uint64 `json:"AmountRequest" ml:"AmountRequest" form:"AmountRequest" validate:"required,minAmount,maxAmount"`
		PartnerCode           string `json:"PartnerCode" ml:"PartnerCode" form:"PartnerCode" validate:"required,partnerExist,partnerActive"`
		PartnerIdentification uint   `json:"PartnerIdentification" ml:"PartnerIdentification" form:"PartnerIdentification" validate:"required"`
		OrderID               uint64 `json:"OrderID" ml:"OrderID" form:"OrderID" validate:"required,orderExist,orderStatusValid"`
		// create order, update amount when partner recharge
		TypeRequest string `json:"TypeRequest" validate:"required,typeRequestExist"`
	}

	RequestBalanceDtoInterface interface {
	}
)

func (a *RequestBalanceDto) BindDataDto(c *fiber.Ctx) (dto_api_response.ResponseRequestBalanceDto, error) {
	if errBP := c.BodyParser(&a); errBP != nil {
		res := dto_api_response.ResponseRequestBalanceDto{
			HttpCode: http_status.StatusBadRequest,
			Status:   dto_api_response.STATUS_ERROR,
			Code:     http_status.StatusBadRequest,
			Message:  "param input not valid, please check doc and try again",
			ErrorDetail: web_server.MapBase{
				"error_list": web_server.MapBase{
					"all_error": errBP.Error(),
				},
			},
		}

		log.Error("param input not valid, please check doc and try again")
		return res, errBP
	}

	return dto_api_response.ResponseRequestBalanceDto{}, nil
}

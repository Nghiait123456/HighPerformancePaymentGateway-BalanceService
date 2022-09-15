package dto_api_request

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_response"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
	log "github.com/sirupsen/logrus"
)

type (
	OneRequestBalanceDto struct {
		AmountRequest         uint64 `json:"AmountRequest" ml:"AmountRequest" form:"AmountRequest" validate:"required,minAmount,maxAmount"`
		PartnerCode           string `json:"PartnerCode" ml:"PartnerCode" form:"PartnerCode" validate:"required,partnerExist,partnerActive"`
		PartnerIdentification uint   `json:"PartnerIdentification" ml:"PartnerIdentification" form:"PartnerIdentification" validate:"required"`
		OrderID               uint64 `json:"OrderID" ml:"OrderID" form:"OrderID" validate:"required,orderExist,orderStatusValid"`
		// create order, update amount when partner recharge
		TypeRequest string `json:"TypeRequest" validate:"required,typeRequestExist"`
	}
	GroupRequestBalanceDto = []OneRequestBalanceDto

	RequestBalanceDto struct {
		Requests GroupRequestBalanceDto
	}

	RequestBalanceDtoInterface interface {
	}
)

func (a *RequestBalanceDto) BindDataDto(c *fiber.Ctx) (dto_api_response.ResponseRequestBalanceDto, error) {
	var g GroupRequestBalanceDto
	if errBP := c.BodyParser(&g); errBP != nil {
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

		errML := fmt.Sprintf("param input not valid, please check doc and try again, detail: %s", errBP.Error())
		log.Error(errML)
		return res, errBP
	}

	a.Requests = g
	return dto_api_response.ResponseRequestBalanceDto{}, nil
}

func NewRequestBalanceDto() *RequestBalanceDto {
	return &RequestBalanceDto{}
}

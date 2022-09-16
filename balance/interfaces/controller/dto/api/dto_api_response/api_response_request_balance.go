package dto_api_response

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
	"github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
)

type (
	ResponseRequestBalanceDto struct {
		HttpCode           int
		Status             string
		Code               int
		Message            string
		ErrorDetail        any //struct json
		ListRequestSuccess any
	}

	ErrorDetailDefault struct {
		ErrorList string //json
	}
)

const (
	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"
)

func (r *ResponseRequestBalanceDto) Response(c web_server.ContextBase) error {
	return c.Status(r.HttpCode).JSON(web_server.MapBase{
		"Status":      r.Status,
		"HttpCode":    r.HttpCode,
		"Message":     r.Message,
		"ErrorDetail": r.ErrorDetail,
	})
}

func (r *ResponseRequestBalanceDto) MappingFrServiceRequestBalanceResponse(listSuccess application.ListRequestSuccess, response respone_request_balance.RequestBalanceResponse, status bool) {
	//todo implement mapping error code to error response
	if status == true {
		r.HttpCode = http_status.StatusOK
		r.Status = STATUS_SUCCESS
		r.Code = response.Code
		r.Message = response.Message
		r.ListRequestSuccess = listSuccess
		r.ErrorDetail = ErrorDetailDefault{}

		return
	}

	if status == false {
		r.HttpCode = http_status.StatusBadRequest
		r.Status = STATUS_ERROR
		r.Code = response.Code
		r.Message = response.Message
		r.ListRequestSuccess = listSuccess
		r.ErrorDetail = ErrorDetailDefault{}
	}
}

func NewResponseRequestBalanceDto() *ResponseRequestBalanceDto {
	return &ResponseRequestBalanceDto{}
}

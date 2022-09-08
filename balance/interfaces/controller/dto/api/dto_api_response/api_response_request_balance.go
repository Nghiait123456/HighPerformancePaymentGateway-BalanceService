package dto_api_response

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/server/web_server"
)

type (
	ResponseRequestBalanceDto struct {
		HttpCode    int
		Status      string
		Code        int
		Message     string
		ErrorDetail web_server.MapBase //json
	}
)

const (
	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"
)

func (r *ResponseRequestBalanceDto) response(c web_server.ContextBase) error {
	return c.Status(r.HttpCode).JSON(web_server.MapBase{
		"Status":      r.Status,
		"HttpCode":    r.HttpCode,
		"Message":     r.Message,
		"ErrorDetail": r.ErrorDetail,
	})
}

package error_http

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/pkg_internal/error_base"

type ErrorInternal = error_base.Error

const (
	ERROR_INTERNAL_CODE    = 500
	ERROR_INTERNAL_MESSAGE = "Internal Error"
)

func NewErrorInternal() ErrorInternal {
	return &error_base.ErrorBase{
		Code:    ERROR_INTERNAL_CODE,
		Message: ERROR_INTERNAL_MESSAGE,
	}
}

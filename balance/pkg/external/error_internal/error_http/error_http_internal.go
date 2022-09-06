package error_http

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorInternal     = error_base.Error
	ErrorInternalType = error_base.ErrorBase
)

const (
	ERROR_INTERNAL_CODE           = 500
	ERROR_INTERNAL_MESSAGE        = "Internal Error"
	ERROR_HTTP_INTERNAL_SIGNATURE = "error_http_ERROR_HTTP_INTERNAL_SIGNATURE"
)

func NewErrorInternal() ErrorInternal {
	return &ErrorInternalType{
		Code:      ERROR_INTERNAL_CODE,
		Message:   ERROR_INTERNAL_MESSAGE,
		Signature: ERROR_HTTP_INTERNAL_SIGNATURE,
	}
}

func IsErrorHttpInternal(e error) bool {
	return error_base.IsErrorOfType(e, ERROR_HTTP_INTERNAL_SIGNATURE)
}

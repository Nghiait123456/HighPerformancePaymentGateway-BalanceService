package error_http

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/pkg_internal/error_base"

type (
	ErrorInternal     = error_base.Error
	ErrorInternalType = error_base.ErrorBase
)

const (
	ERROR_HTTP_INTERNAL_SIGNATURE = "ERROR_HTTP_INTERNAL_SIGNATURE"
	ERROR_INTERNAL_CODE           = 500
	ERROR_INTERNAL_MESSAGE        = "Internal Error"
)

func NewErrorInternal() ErrorInternal {
	return &ErrorInternalType{
		Code:      ERROR_INTERNAL_CODE,
		Message:   ERROR_INTERNAL_MESSAGE,
		Signature: ERROR_HTTP_INTERNAL_SIGNATURE,
	}
}

func IsErrorHttpInternal(e ErrorInternal) bool {
	return e.GetSignature() == ERROR_HTTP_INTERNAL_SIGNATURE
}

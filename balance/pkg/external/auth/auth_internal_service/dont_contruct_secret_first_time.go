package auth_internal_service

import "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/error_base"

type (
	ErrorDontConstructSecretFirstTime          = error_base.ErrorBase
	ErrorDontConstructSecretFirstTimeInterface = error_base.Error
)

const (
	ERROR_CODE_DONT_CONTRUCT_SECRET_FIRST_TIME      = 600
	ERROR_MESSAGE_DONT_CONTRUCT_SECRET_FIRST_TIME   = "Dont construct secret in first time"
	ERROR_SIGNATURE_DONT_CONTRUCT_SECRET_FIRST_TIME = "auth_internal_service_ERROR_DONT_CONTRUCT_SECRET_FIRST_TIME"
)

func NewErrorDontConstructSecretFirstTime() ErrorDontConstructSecretFirstTimeInterface {
	return &ErrorDontConstructSecretFirstTime{
		Code:      ERROR_CODE_DONT_CONTRUCT_SECRET_FIRST_TIME,
		Message:   ERROR_MESSAGE_DONT_CONTRUCT_SECRET_FIRST_TIME,
		Signature: ERROR_SIGNATURE_DONT_CONTRUCT_SECRET_FIRST_TIME,
	}
}

func IsErrorDontConstructSecretFirstTime(e error) bool {
	if error_base.IsErrorBase(e) != true {
		return false
	}

	err := error_base.GetErrorBase(e)
	return err.GetSignature() == ERROR_SIGNATURE_DONT_CONTRUCT_SECRET_FIRST_TIME
}

package respone_request_balance

const (
	ERROR_CODE_PARTNERCODE_DOES_NOT_EXITS    = 404
	ERROR_MESSAGE_PARTNERCODE_DOES_NOT_EXITS = "partner does not exist"
)

func ErrPartnerCodeDoesNotExist() RequestBalanceResponse {
	return RequestBalanceResponse{
		Status:  STATUS_ERROR,
		Code:    ERROR_CODE_PARTNERCODE_DOES_NOT_EXITS,
		Message: ERROR_MESSAGE_PARTNERCODE_DOES_NOT_EXITS,
	}
}

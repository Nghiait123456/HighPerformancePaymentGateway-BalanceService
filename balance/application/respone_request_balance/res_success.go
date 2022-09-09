package respone_request_balance

const (
	CODE_SUCCESS    = 0
	MESSAGE_SUCCESS = "request balance success, but balance is processing, please check other event response"
)

func SuccessBalanceResponse() RequestBalanceResponse {
	return RequestBalanceResponse{
		Status:  STATUS_ERROR,
		Code:    CODE_SUCCESS,
		Message: MESSAGE_SUCCESS,
	}
}

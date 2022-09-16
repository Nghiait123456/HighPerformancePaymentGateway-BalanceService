package respone_request_balance

const (
	CODE_ERROR_ESTOP    = 2
	MESSAGE_ERROR_ESTOP = "have some error, systems is stopping, please try again late"
)

func ErrorWhySystemEStop() RequestBalanceResponse {
	return RequestBalanceResponse{
		Status:  STATUS_ERROR,
		Code:    CODE_ERROR_ESTOP,
		Message: MESSAGE_ERROR_ESTOP,
	}
}

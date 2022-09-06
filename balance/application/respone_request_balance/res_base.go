package respone_request_balance

type (
	RequestBalanceResponse struct {
		Status  string
		Code    int
		Message string
	}
)

const (
	ERR_STATUS     = "error"
	SUCCESS_STATUS = "success"
)

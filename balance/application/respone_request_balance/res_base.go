package respone_request_balance

type (
	RequestBalanceResponse struct {
		Status  string
		Code    int
		Message string
	}
)

const (
	STATUS_ERROR   = "error"
	STATUS_SUCCESS = "success"
)

func (r RequestBalanceResponse) IsError() bool {
	return r.Status == STATUS_ERROR
}

func (r RequestBalanceResponse) IsSuccess() bool {
	return r.Status == STATUS_SUCCESS
}

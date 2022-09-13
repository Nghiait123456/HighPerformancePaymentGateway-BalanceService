package handle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
	"github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"
	validate_api "github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/api/validate"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_request"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_response"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
	validate_base "github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/validate"
)

type (
	RequestBalance struct {
		sv application.ServiceInterface
	}

	RequestBalanceResponse struct {
		HttpStatus int
		Status     string
		Code       int
		Message    string
	}
)

func (r *RequestBalance) Inject(
	sv application.ServiceInterface,
) {
	r.sv = sv
}

func (r *RequestBalance) HandleOneRequestBalance(c *fiber.Ctx) error {
	rqDto := dto_api_request.NewRequestBalanceDto()
	res, errB := rqDto.BindDataDto(c)
	if errB != nil {
		return res.Response(c)
	}

	validate := validate_api.ValidateApiRequestBalance{
		VB:  validate_base.NewBaseValidate(),
		Dto: *rqDto,
	}

	validate.Init()
	resV, errV := validate.Validate()
	if errV != nil {
		return resV.Response(c)
	}

	bRequest := application.BalanceRequest{
		BRequest: calculator.BalancerRequest{
			AmountRequest:         rqDto.AmountRequest,
			PartnerCode:           rqDto.PartnerCode,
			PartnerIdentification: rqDto.PartnerIdentification,
			OrderID:               rqDto.OrderID,
			TypeRequest:           rqDto.TypeRequest,
		},
	}
	rs, ErrorHRB := r.sv.HandleRequestBalance(bRequest)

	resProcess := dto_api_response.NewResponseRequestBalanceDto()
	resProcess.MappingFrServiceRequestBalanceResponse(rs, ErrorHRB)
	return resProcess.Response(c)
}

func (r *RequestBalance) TransferToResponseHttp(resRB respone_request_balance.RequestBalanceResponse) RequestBalanceResponse {
	switch resRB.Status {
	case respone_request_balance.STATUS_SUCCESS:
		return r.ResponseSuccess(resRB)
	default:
		return r.ResponseErrorDefault(resRB)
	}
}

func (r RequestBalance) ResponseSuccess(resRB respone_request_balance.RequestBalanceResponse) RequestBalanceResponse {
	return RequestBalanceResponse{
		HttpStatus: http_status.StatusOK,
		Status:     resRB.Status,
		Code:       resRB.Code,
		Message:    resRB.Message,
	}
}

func (r RequestBalance) ResponseErrorDefault(resRB respone_request_balance.RequestBalanceResponse) RequestBalanceResponse {
	return RequestBalanceResponse{
		HttpStatus: http_status.StatusInternalServerError,
		Status:     resRB.Status,
		Code:       resRB.Code,
		Message:    resRB.Message,
	}
}

func NewRequestBalance(sv application.ServiceInterface) *RequestBalance {
	return &RequestBalance{
		sv: sv,
	}
}

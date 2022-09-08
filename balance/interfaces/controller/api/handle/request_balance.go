package handle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
	"github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/interfaces/controller/dto/api/dto_api_request"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/http/http_status"
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
	rq := []dto_api_request.RequestBalanceDto{}
	if err := c.BodyParser(&rq); err != nil {
		fmt.Println("parser Error")
		return err
	}

	log.Printf("%#v\n", persons)
	// []main.Person{main.Person{Name:"john", Pass:"doe"}, main.Person{Name:"jack", Pass:"jill"}}
	return c.SendString("Post Called")
	//https://stackoverflow.com/questions/71934014/how-to-parse-post-request-body-with-arbitrary-number-of-parameters-in-go-fiber
}

func (r *RequestBalance) TransferToResponseHttp(resRB respone_request_balance.RequestBalanceResponse) RequestBalanceResponse {
	switch resRB.Status {
	case respone_request_balance.SUCCESS_STATUS:
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

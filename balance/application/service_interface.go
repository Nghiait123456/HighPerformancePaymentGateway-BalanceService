package application

import "github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"

type (
	ServiceInterface interface {
		Init()
		HandleRequestBalance(b BalanceRequest) (bool, respone_request_balance.RequestBalanceResponse)
	}
)

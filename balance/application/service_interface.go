package application

import "github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"

type (
	ServiceInterface interface {
		Init()
		HandleOneRequestBalance(b BalanceRequest) (respone_request_balance.RequestBalanceResponse, bool)
		HandleGroupRequestBalance(gb GroupBalanceRequest) (respone_request_balance.RequestBalanceResponse, DetailResultGroupRequest, bool)
	}
)

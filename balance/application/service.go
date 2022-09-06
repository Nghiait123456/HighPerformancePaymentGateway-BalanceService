package application

import "github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"

type (
	Service struct {
		allP AllPartnerBalanceInterface
	}
)

func (s *Service) Init() {
	s.allP.Init()
}

func (s *Service) HandleRequestBalance(b BalanceRequest) (bool, respone_request_balance.RequestBalanceResponse) {
	return s.allP.HandleRequestBalance(b)
}

func NewService() ServiceInterface {
	return &Service{}
}

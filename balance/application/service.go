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

func (s *Service) HandleOneRequestBalance(b BalanceRequest) (respone_request_balance.RequestBalanceResponse, bool) {
	return s.allP.HandleOneRequestBalance(b)
}

func (s *Service) HandleGroupRequestBalance(gb GroupBalanceRequest) (respone_request_balance.RequestBalanceResponse, DetailResultGroupRequest, bool) {
	return s.allP.HandleGroupRequestBalance(gb)
}

func NewService(allP AllPartnerBalanceInterface) *Service {
	var _ ServiceInterface = (*Service)(nil)
	return &Service{allP: allP}
}

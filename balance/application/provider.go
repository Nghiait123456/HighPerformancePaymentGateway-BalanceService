package application

import (
	"github.com/google/wire"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/queue_job_request"
)

var ProviderAllPartnerBalance = wire.NewSet(
	NewAllPartnerBalance,
	calculator.ProviderAllPartner,
	queue_job_request.NewQueueJob,
	wire.Bind(new(AllPartnerBalanceInterface), new(*AllPartnerBalance)),
)

var ProviderService = wire.NewSet(
	NewService, ProviderAllPartnerBalance, wire.Bind(new(ServiceInterface), new(*Service)),
)

package handle

import (
	"github.com/google/wire"
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
)

var ProviderRequestBalance = wire.NewSet(
	NewRequestBalance,
	application.ProviderService,
)

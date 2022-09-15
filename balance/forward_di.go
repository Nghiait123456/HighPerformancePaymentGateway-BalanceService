package balance

import "github.com/high-performance-payment-gateway/balance-service/balance/application"

/**
forward to wire_gen
*/

func ForwardProviderService() application.ServiceInterface {
	return ProviderService()
}

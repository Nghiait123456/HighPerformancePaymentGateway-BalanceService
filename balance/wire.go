//go:build wireinject
// +build wireinject

package balance

import (
	"github.com/google/wire"
	"github.com/high-performance-payment-gateway/balance-service/balance/application"
)

func ProviderService() application.ServiceInterface {
	wire.Build(
		application.ProviderService,
	)

	return &application.Service{}
}

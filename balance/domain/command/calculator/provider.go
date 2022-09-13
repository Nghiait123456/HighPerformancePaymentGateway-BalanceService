package calculator

import (
	"github.com/google/wire"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/logs_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
)

func NewPartnersBalance() PartnersBalance {
	return make(PartnersBalance)
}

func InstanceCnRechargeLogs() sql.Connect {
	// todo get from globle config
	var temp sql.Connect
	return temp
}

func InstanceCnBalance() sql.Connect {
	// todo get from globle config
	var temp sql.Connect
	return temp
}

var ProviderAllPartner = wire.NewSet(
	NewAllPartner,
	NewPartnersBalance,
	InstanceCnRechargeLogs,
	InstanceCnBalance(),
	logs_request_balance.NewLog,
	wire.Bind(new(AllPartnerInterface), new(*AllPartner)),
)

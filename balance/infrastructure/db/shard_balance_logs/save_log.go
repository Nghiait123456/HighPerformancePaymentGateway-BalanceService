package shard_balance_logs

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
)

type (
	SaveLogRequestBalance struct {
	}

	SaveLogRequestBalanceInterface interface {
		Save(lb LBShardLogInterface, log orm.BalanceRequestLog) error
	}
)

func (s SaveLogRequestBalance) Save(lb LBShardLogInterface, log orm.BalanceRequestLog) error {
	rp := repository.NewBalanceLogRepository()
	cn := lb.loadBalanceShard(log.PartnerCode)
	rp.SetConnect(cn)

	ok := rp.CreateNewWithTrans(log)
	if ok != nil {
		return ok
	}

	return nil
}

func NewSaveLog() SaveLogRequestBalanceInterface {
	return &SaveLogRequestBalance{}
}

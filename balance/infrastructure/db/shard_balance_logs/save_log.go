package shard_balance_logs

import (
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
)

type (
	SaveLog struct {
	}

	SaveLogInterface interface {
		Save(lb LBShardLogInterface, log orm.BalanceLog) error
	}
)

func (s SaveLog) Save(lb LBShardLogInterface, log orm.BalanceLog) error {
	rp := repository.NewBalanceLogRepository()
	cn := lb.loadBalanceShard(log.PartnerCode)
	rp.SetConnect(cn)

	ok := rp.CreateNewWithTrans(log)
	if ok != nil {
		return ok
	}

	return nil
}

func NewSaveLog() SaveLogInterface {
	return &SaveLog{}
}

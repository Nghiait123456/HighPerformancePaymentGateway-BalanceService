package entity

import "github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"

type PartnerBalanceShards []orm.PartnerBalanceShard

func (e PartnerBalanceShards) GetAllShardCode() []string {
	var rs []string

	for _, v := range e {
		rs = append(rs, v.ShardCode)
	}

	return rs
}

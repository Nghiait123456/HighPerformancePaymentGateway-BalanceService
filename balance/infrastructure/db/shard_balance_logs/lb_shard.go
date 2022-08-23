package shard_balance_logs

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository_other_service"
	"math"
	"os"
	"sync/atomic"
)

/**
lb for infra sharding table balance_log
*/
type LBShardLog struct {
	allPShard  allPShard
	allConnect map[string]Connect //[shardCode]Connect

}

type LBShardLogInterface interface {
	// loadBalanceShard select one shard in pool
	loadBalanceShard(partnerCode string) Connect
	InitAllShard() error
	UpdateConnectAllShard() error
}

func (lb *LBShardLog) loadBalanceShard(partnerCode string) Connect {
	partnerShard, ok := lb.allPShard[partnerCode]
	if !ok {

	}
	if partnerShard.totalShard == 1 {
		for _, v := range partnerShard.listConnect {
			return v
		}
	}

	//Round Robin
	atomic.AddUint64(&partnerShard.indexShard, 1)
	if partnerShard.indexShard == math.MaxUint64 {
		partnerShard.indexShard = 0
	}

	shardId := partnerShard.indexShard % uint64(partnerShard.totalShard)
	if shardId > uint64(len(partnerShard.listCodeShard))-1 {
		panic(fmt.Sprintf("totalShard not match listCodeShard, partnerCode %s", partnerCode))
		os.Exit(0)
	}
	shardCode := partnerShard.listCodeShard[uint(shardId)]

	return partnerShard.listConnect[shardCode]
}

func (lb *LBShardLog) InitAllShard() error {
	banlancerShardRp := repository.NewBalanceShardRepository()
	partnerBalanceShardRp := repository.NewPartnerBalanceShardRepository()
	parnterRp := repository_other_service.NewPartnerRepository()

	// init connect all shard
	allBalanceShardAct, errABl := banlancerShardRp.AllBalanceShardActive()
	if errABl != nil {
		return errABl
	}

	for _, v := range allBalanceShardAct {
		db, err := NewConnectShard().ConnectOneShard(v.GetShardDsn())
		if err != nil {
			errMessage := fmt.Sprintf("don't connect DB shard balance why error: %s", err.Error())
			panic(errMessage)
			os.Exit(0)
		}

		lb.allConnect[v.ShardCode] = db
	}

	// map connect to partner
	allPartnerCodeAct := parnterRp.AllPartnerCodeActive()

	for _, partnerCode := range allPartnerCodeAct {
		var oneP onePartnerShard
		oneBalanceShard, ok := allBalanceShardAct[partnerCode]
		if !ok {
			errM := fmt.Sprintf("missing config balance shard of parterCode: %s", partnerCode)
			panic(errM)
			os.Exit(0)
		}

		allShardOfPartner, errShardOP := partnerBalanceShardRp.GetAllActiveByPartner(partnerCode)
		if errShardOP != nil {
			return errShardOP
		}

		allShardCode := allShardOfPartner.GetAllShardCode()

		oneP.partnerCode = partnerCode
		oneP.status = "active"
		oneP.totalShard = uint(len(allShardOfPartner))
		oneP.indexShard = 0
		oneP.listCodeShard = allShardCode

		//init all shard one partner
		for _, v := range allShardOfPartner {
			var oneShard Shard
			oneShard.shardId = v.ShardId
			oneShard.shardCode = v.ShardCode
			oneShard.dsnEncry = oneBalanceShard.ShardDsnEncry
			oneShard.dnsRaw = oneBalanceShard.GetShardDsn()
			oneShard.db, ok = lb.allConnect[v.ShardCode]
			if !ok {
				panic(fmt.Sprintf("shardCode %s don't match in allconect %v", v.ShardCode, lb.allConnect))
			}

			oneP.listConnect[v.ShardCode] = oneShard.db
			oneP.listShard[v.ShardCode] = oneShard
			oneP.trafficShard[v.ShardCode] = trafficShard{}
		}

		lb.allPShard[partnerCode] = oneP
	}

	return nil
}

func (lb *LBShardLog) UpdateConnectAllShard() error {
	return lb.InitAllShard()
}

func NewLBShardLog() LBShardLogInterface {
	return &LBShardLog{}
}

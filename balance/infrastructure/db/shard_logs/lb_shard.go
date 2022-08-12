package shard_logs

import (
	"context"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
	"math"
	"sync/atomic"
	"time"
)

type lbShard struct {
	allPShard allPShard
}

type lbShardInterface interface {
	// loadBalanceShard select one shard in pool
	loadBalanceShard(partnerCode string) Shard
	InitConnectAllShard() error
	UpdateConnectAllShard() error
}

func (lb *lbShard) loadBalanceShard(partnerCode string) Shard {
	partnerShard, ok := lb.allPShard[partnerCode]
	if !ok {

	}
	if partnerShard.totalShard == 1 {
		for _, v := range partnerShard.listShard {
			return v
		}
	}

	//Round Robin
	atomic.AddUint64(&partnerShard.indexShard, 1)
	if partnerShard.indexShard == math.MaxUint64 {
		partnerShard.indexShard = 0
	}
	shardId := partnerShard.indexShard % uint64(partnerShard.totalShard)

	return partnerShard.listShard[uint(shardId)]
}

func (lb *lbShard) InitAllShard() error {
	// get all ShardBalane

	bsr := repository.NewBalanceShardRepository(ctx)

	//get all Partner
	ctxP, cancelP := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancelP()
	pbs := repository.NewPartnerBalanceShardRepository(ctxP)

	allBsr, err := bsr.AllBalanceShard()
	if err != nil {
		return err
	}

	//allShardBalance :=
}

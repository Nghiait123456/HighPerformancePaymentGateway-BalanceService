package shard_logs

import "sync/atomic"

type lbShard struct {
	allPShard allPShard
}

type lbShardInterface interface {
	// loadBalanceShard select one shard in pool
	loadBalanceShard(partnerCode string) Shard
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

	atomic.AddUint64(&partnerShard.indexShard, 1)
	//Round Robin
	shardId := partnerShard.indexShard % uint64(partnerShard.totalShard)
	return partnerShard.listShard[uint(shardId)]
}

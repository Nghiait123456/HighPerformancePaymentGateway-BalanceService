package balance_shard_logs

type lbShard struct {
	partnerShard onePartnerShard
}

type lbShardInterface interface {
	// loadBalanceShard select one shard in pool
	loadBalanceShard() Shard
}

func (lb *lbShard) loadBalanceShard() Shard {
	if lb.partnerShard.totalShard == 1 {
		for _, v := range lb.partnerShard.listShard {
			return v
		}
	}

	//Round Robin
	shardId := lb.partnerShard.indexShard % uint64(lb.partnerShard.totalShard)
	return lb.partnerShard.listShard[uint(shardId)]
}

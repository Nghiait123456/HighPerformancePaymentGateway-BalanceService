package balance_shard_logs

/**
  all info partner config sharding
*/
type allPartnerShard struct {
}

type onePartnerShard struct {
	partnerCode  string
	status       string
	totalShard   uint
	indexShard   uint64
	listShard    map[uint]Shard
	trafficShard map[uint]trafficShard
}

type trafficShard struct {
	totalVisit uint32
}

type Shard struct {
	shardId   uint
	shardCode string
}

type allPShard map[string]onePartnerShard

type partnerShardingInterface interface {
	allConfigPartnerSharding() allPShard
	getAllShardOnePartner(partnerCode string)
}

package shard_logs

/**
  all info partner config sharding
*/
type allPartnerShard struct {
}

//onePartnerShard partner mapping shard
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

// Shard : infor one shard
type Shard struct {
	shardId   uint
	shardCode string
}

type allPShard map[string]onePartnerShard

type partnerShardingInterface interface {
	allConfigPartnerSharding() allPShard
	getAllShardOnePartner(partnerCode string)
}

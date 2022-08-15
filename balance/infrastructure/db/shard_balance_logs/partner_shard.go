package shard_balance_logs

//onePartnerShard partner mapping shard
type onePartnerShard struct {
	partnerCode   string
	status        string
	totalShard    uint
	indexShard    uint64
	listCodeShard []string                // [1]ShardCode1, [2]ShardCode2
	listShard     map[string]Shard        // [shardCode]Shard
	listConnect   map[string]Connect      // [shardCode]*gorm.DB
	trafficShard  map[string]trafficShard // [shardCode]trafficShard
}

type trafficShard struct {
	totalVisit uint64
}

// Shard : info one shard
type Shard struct {
	shardId   uint32
	shardCode string
	dsnEncry  string
	dnsRaw    string
	db        Connect
}

type allPShard map[string]onePartnerShard //[partnerCode]onePartnerShard

type partnerShardingInterface interface {
	allConfigPartnerSharding() allPShard
	getAllShardOnePartner(partnerCode string)
}

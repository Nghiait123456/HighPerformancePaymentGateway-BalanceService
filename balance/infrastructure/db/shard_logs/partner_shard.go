package shard_logs

import "gorm.io/gorm"

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

// Shard : info one shard
type Shard struct {
	shardId   uint
	shardCode string
	dsnEncry  string
	dnsRaw    string
	db        *gorm.DB
}

type allPShard map[string]onePartnerShard

type partnerShardingInterface interface {
	allConfigPartnerSharding() allPShard
	getAllShardOnePartner(partnerCode string)
}

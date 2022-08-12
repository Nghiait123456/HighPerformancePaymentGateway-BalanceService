package orm

type BalanceShard struct {
	ID            uint32 `gorm:"<-:create"`
	ShardName     string
	ShardCode     string
	ShardDnsEncry string
	Status        string
	CreatedAt     uint32
	UpdatedAt     uint32
}

// TableName overrides
func (b *BalanceShard) TableName() string {
	return "balance_shard"
}

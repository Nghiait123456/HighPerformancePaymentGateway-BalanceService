package orm

type BalanceShard struct {
	ID            uint32 `gorm:"<-:create"`
	ShardName     string
	ShardCode     string
	ShardDsnEncry string
	Status        string
	CreatedAt     uint32
	UpdatedAt     uint32
}

const (
	BALANCE_SHARD_STATUS_ACTIVE = "active"
)

// TableName overrides
func (b *BalanceShard) TableName() string {
	return "balance_shard"
}

func (b BalanceShard) IsActive() bool {
	return b.Status == BALANCE_SHARD_STATUS_ACTIVE
}

func (b BalanceShard) StatusActive() string {
	return BALANCE_SHARD_STATUS_ACTIVE
}

func (b BalanceShard) GetShardDsn() string {
	// todo add describe
	return b.ShardDsnEncry
}

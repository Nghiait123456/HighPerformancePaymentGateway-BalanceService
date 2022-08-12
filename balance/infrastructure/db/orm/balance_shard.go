package orm

type BalanceShard struct {
	ID          uint32 `gorm:"<-:create"`
	PartnerCode string
	ShardName   string
	ShardCode   string
	ShardLink   string
	Status      string
	CreatedAt   uint32
	UpdatedAt   uint32
}

// TableName overrides
func (b *BalanceShard) TableName() string {
	return "balance_shard"
}

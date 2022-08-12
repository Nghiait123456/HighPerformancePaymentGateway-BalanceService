package orm

type PartnerBalanceShard struct {
	ID          uint32 `gorm:"<-:create"`
	PartnerCode string
	ShardId     uint32
	Status      string
	CreateAt    uint32
	UpdatedAt   uint32
}

// TableName overrides
func (b *PartnerBalanceShard) TableName() string {
	return "partner_balance_shard"
}

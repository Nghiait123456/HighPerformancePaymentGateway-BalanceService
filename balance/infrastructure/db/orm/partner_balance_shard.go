package orm

type PartnerBalanceShard struct {
	ID          uint32 `gorm:"<-:create"`
	PartnerCode string
	ShardId     uint32
	ShardCode   string
	Status      string
	CreateAt    uint32
	UpdatedAt   uint32
}

const (
	PARTNER_BALANCE_SHARD_STATUS_ACTIVE = "active"
)

// TableName overrides
func (b *PartnerBalanceShard) TableName() string {
	return "partner_balance_shard"
}

func (b PartnerBalanceShard) IsActive() bool {
	return b.Status == PARTNER_BALANCE_SHARD_STATUS_ACTIVE
}

func (b PartnerBalanceShard) StatusActive() string {
	return PARTNER_BALANCE_SHARD_STATUS_ACTIVE
}

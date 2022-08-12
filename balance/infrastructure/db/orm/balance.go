package orm

type Balance struct {
	ID          uint32 `gorm:"<-:create"`
	Balance     uint32
	PartnerCode string
	Status      string
	CreatedAt   uint32
	UpdatedAt   uint32
}

const (
	STATUS_ACTIVE = "active"
)

// TableName overrides
func (b *Balance) TableName() string {
	return "balances"
}

func (b *Balance) IsActive() bool {
	return b.Status == STATUS_ACTIVE
}

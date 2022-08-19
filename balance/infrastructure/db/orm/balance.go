package orm

type Balance struct {
	ID                    uint32 `gorm:"<-:create"`
	Balance               uint64
	PartnerCode           string
	Status                string
	IndexLogRequestLatest uint64
	CreatedAt             uint32
	UpdatedAt             uint32
}

const (
	BALANCE_STATUS_ACTIVE = "active"
)

// TableName overrides
func (b *Balance) TableName() string {
	return "balances"
}

func (b *Balance) IsActive() bool {
	return b.Status == BALANCE_STATUS_ACTIVE
}

func (b *Balance) StatusActive() string {
	return BALANCE_STATUS_ACTIVE
}

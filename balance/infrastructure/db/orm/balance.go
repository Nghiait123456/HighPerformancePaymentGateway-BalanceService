package orm

type Balance struct {
	ID          uint32 `gorm:"primaryKey;<-:create"`
	Balance     uint32
	PartnerCode string
	Status      string
	CreatedAt   uint32
	UpdatedAt   uint32
}

const (
	STATUS_ACTIVE = "active"
)

func (b *Balance) IsActive() bool {
	return b.Status == STATUS_ACTIVE
}

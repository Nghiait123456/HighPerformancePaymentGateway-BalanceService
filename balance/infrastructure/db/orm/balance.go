package orm

type Balance struct {
	ID          uint32 `gorm:"primaryKey;<-:create"`
	Balance     uint32
	PartnerCode string
	status      string
	CreatedAt   uint32
	UpdatedAt   uint32
}

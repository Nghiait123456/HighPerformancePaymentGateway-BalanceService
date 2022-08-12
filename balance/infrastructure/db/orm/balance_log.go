package orm

type BalanceLog struct {
	Id        uint32 `gorm:"primaryKey;<-:create"`
	OrderId   uint64 `gorm:"uniqueIndex"`
	Amount    uint32
	Balance   uint32
	Status    string
	CreatedAt uint32
	UpdatedAt uint32
}

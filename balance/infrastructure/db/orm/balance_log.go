package orm

type BalanceLog struct {
	ID        uint32 `gorm:"<-:create"`
	OrderId   uint64 `gorm:"uniqueIndex"`
	Amount    uint32
	Balance   uint32
	Status    string
	CreatedAt uint32
	UpdatedAt uint32
}

// TableName overrides
func (b *BalanceLog) TableName() string {
	return "balances_logs"
}
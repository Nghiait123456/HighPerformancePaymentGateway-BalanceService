package orm

type BalanceLog struct {
	ID                uint32 `gorm:"<-:create"`
	OrderId           uint64 `gorm:"uniqueIndex"`
	PartnerCode       string
	AmountRequest     uint64
	AmountPlaceHolder uint64
	Balance           uint64
	Status            string
	CreatedAt         uint32
	UpdatedAt         uint32
}

const (
	BALANCE_LOG_STATUS_PROCESSING = "processing"
)

// TableName overrides
func (b *BalanceLog) TableName() string {
	return "balances_logs"
}

func (b BalanceLog) StatusProcessing() string {
	return BALANCE_LOG_STATUS_PROCESSING
}

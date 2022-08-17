package orm

type BalanceRequestLog struct {
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
func (b *BalanceRequestLog) TableName() string {
	return "balances_logs"
}

func (b BalanceRequestLog) StatusProcessing() string {
	return BALANCE_LOG_STATUS_PROCESSING
}

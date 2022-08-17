package orm

type BalanceRechargeLog struct {
	ID             uint32 `gorm:"<-:create"`
	OrderId        uint64 `gorm:"uniqueIndex"`
	PartnerCode    string
	AmountRecharge uint64
	Balance        uint64
	Status         string
	CreatedAt      uint32
	UpdatedAt      uint32
}

const (
	BALANCE_RECHARGE_LOG_STATUS_SUCCESS = "success"
)

// TableName overrides
func (b *BalanceRechargeLog) TableName() string {
	return "balances_logs"
}

func (b BalanceRechargeLog) StatusSuccess() string {
	return BALANCE_RECHARGE_LOG_STATUS_SUCCESS
}

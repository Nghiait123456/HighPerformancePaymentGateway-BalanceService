package calculator

type (
	balancePlaceHolderHistory struct {
		allAmountPlaceHolder allAmountPlaceHolderFrLogs
	}

	balancePlaceHolderHistoryInterface interface {
		allPartnerCode() []string
		loadAllPlaceHolderAmountFrLogs() allAmountPlaceHolderFrLogs
		initBalancePlaceHolderHistory() error
		GetAllPlaceHolder() allAmountPlaceHolderFrLogs
	}

	//amountPlaceHolderFrLogs
	amountPlaceHolderFrLogs struct {
		partnerCode       string
		amountPlaceHolder uint64
	}

	allAmountPlaceHolderFrLogs map[string]amountPlaceHolderFrLogs // [PartnerCode]amountPlaceHolderFrLogs
)

func (b *balancePlaceHolderHistory) allPartnerCode() []string {
	fake := []string{"test"}
	return fake
}

func (b *balancePlaceHolderHistory) loadAllPlaceHolderAmountFrLogs() allAmountPlaceHolderFrLogs {
	fake := allAmountPlaceHolderFrLogs{}
	fake["test"] = amountPlaceHolderFrLogs{
		partnerCode:       "test",
		amountPlaceHolder: 50000,
	}
	return fake
}

func (b *balancePlaceHolderHistory) initBalancePlaceHolderHistory() error {
	allAmountPlaceHolder := b.loadAllPlaceHolderAmountFrLogs()
	for k, v := range allAmountPlaceHolder {
		b.allAmountPlaceHolder[k] = v
	}

	return nil
}

func (b *balancePlaceHolderHistory) GetAllPlaceHolder() allAmountPlaceHolderFrLogs {
	return b.allAmountPlaceHolder
}

func NewBalancePlaceHolderHistory() balancePlaceHolderHistoryInterface {
	b := balancePlaceHolderHistory{
		allAmountPlaceHolder: make(allAmountPlaceHolderFrLogs),
	}
	b.initBalancePlaceHolderHistory()
	return &b
}

package calculator

type balancePlaceHolderHistory struct {
	allAmountPlaceHolder allAmountPlaceHolderFrLogs
}

//amountPlaceHolderFrLogs
type amountPlaceHolderFrLogs struct {
	partnerCode       string
	amountPlaceHolder uint
}

type allAmountPlaceHolderFrLogs map[string]amountPlaceHolderFrLogs

type balancePlaceHolderHistoryInterface interface {
	allPartnerCode() []string
	loadAllPlaceHolderAmountFrLogs() allAmountPlaceHolderFrLogs
	initBalancePlaceHolderHistory() error
	GetAllPlaceHolder() allAmountPlaceHolderFrLogs
}

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
	b := balancePlaceHolderHistory{}
	b.initBalancePlaceHolderHistory()
	return &b
}

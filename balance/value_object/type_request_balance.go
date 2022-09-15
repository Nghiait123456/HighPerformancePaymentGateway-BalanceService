package value_object

import "golang.org/x/exp/slices"

const (
	TYPE_REQUEST_BALANCE_REQUEST_PAYMENT  = "payment"
	TYPE_REQUEST_BALANCE_REQUEST_RECHARGE = "recharge"
)

type (
	TypeRequestBalance struct {
		RequestPayment  string
		RequestRecharge string
		AllType         []string //[value][value]
	}
)

func (t *TypeRequestBalance) TypeRequestExist(typeRq string) bool {
	return slices.Contains(t.AllType, typeRq)
}

func (t *TypeRequestBalance) SetAllType() {
	t.AllType = append(t.AllType, TYPE_REQUEST_BALANCE_REQUEST_PAYMENT)
	t.AllType = append(t.AllType, TYPE_REQUEST_BALANCE_REQUEST_RECHARGE)
}

func (t *TypeRequestBalance) Init() {
	t.RequestPayment = TYPE_REQUEST_BALANCE_REQUEST_PAYMENT
	t.RequestRecharge = TYPE_REQUEST_BALANCE_REQUEST_RECHARGE
	t.SetAllType()
}

func NewTypeRequestBalance() *TypeRequestBalance {
	t := TypeRequestBalance{}
	t.Init()
	return &t
}

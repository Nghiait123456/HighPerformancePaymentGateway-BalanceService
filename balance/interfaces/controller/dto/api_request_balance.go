package dto

type (
	ApiRequestBalanceDto struct {
		AmountRequest         uint64 `json:"AmountRequest"`
		PartnerCode           string `json:"PartnerCode"`
		PartnerIdentification uint   `json:"PartnerIdentification"`
		OrderID               uint64 `json:"OrderID"`
		// create order, update amount when partner recharge
		TypeRequest string `json:"TypeRequest"`
	}
)

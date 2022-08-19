package calculator

import (
	"errors"
	"fmt"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/shard_balance_logs"
	"sync"
)

/**
all information partner for calculator balance
*/

type (
	allPartner struct {
		allPartner        map[string]partnerBalance
		muLock            sync.Mutex
		lbShardBalanceLog shard_balance_logs.LBShardLogInterface
		cnRechargeLog     sql.Connect
		cnBalance         sql.Connect
	}

	allPartnerInterface interface {
		initPartnersInterface
	}

	initPartnersInterface interface {
		LoadAllPartnerInfo() (map[string]partnerBalance, error)
		InitAllPartnerInfo() error
		UpdateOnePartner(p partnerBalance) error
		GetOnePartner(partnerCode string) (partnerBalance, error)
		getKeyOnePartner(p partnerBalance) string
	}
)

func (allP *allPartner) LoadAllPartnerInfo() (map[string]partnerBalance, error) {
	fake := make(map[string]partnerBalance)
	//todo get indexLogRequestLatest from DB and update to
	fake["partner_test"] = partnerBalance{
		partnerCode:           "partner_test",
		partnerName:           "TEST",
		partnerIdentification: 1,
		balance:               1000000,
		amountPlaceHolder:     4,
		status:                "active",
		muLock:                sync.Mutex{},
		lbShardBalanceLog:     allP.lbShardBalanceLog,
		cnRechargeLog:         allP.cnRechargeLog,
		cnBalance:             allP.cnBalance,
	}

	return fake, nil
}

func (allP *allPartner) InitAllPartnerInfo() error {
	allPartner, err := allP.LoadAllPartnerInfo()
	if err != nil {
		return err
	}

	balancePlaceHolderHistory := NewBalancePlaceHolderHistory()

	allP.muLock.Lock()
	//init raw allPartner
	for k, v := range allPartner {
		allP.allPartner[k] = v
	}

	//merger balancePlaceHolderHistory to allPartner
	for k, v := range allP.allPartner {
		placeHolder, ok := balancePlaceHolderHistory.GetAllPlaceHolder()[k]
		if ok {
			v.amountPlaceHolder = placeHolder.amountPlaceHolder
			allP.allPartner[k] = v
		} else {
			v.amountPlaceHolder = 0
		}
	}

	allP.muLock.Unlock()

	return nil
}

func (allP *allPartner) UpdateOnePartner(p partnerBalance) error {
	p.muLock.Lock()
	key := allP.getKeyOnePartner(p)
	allP.allPartner[key] = p
	return nil
}

func (allP *allPartner) getKeyOnePartner(p partnerBalance) string {
	return p.partnerCode
}

func (allP *allPartner) GetOnePartner(partnerCode string) (partnerBalance, error) {
	partner, ok := allP.allPartner[partnerCode]
	if !ok {
		err := fmt.Sprintf("partnercode %s not exists", partnerCode)
		return partnerBalance{}, errors.New(err)
	}

	return partner, nil
}

func InitAllPartnerData() allPartnerInterface {
	allPartner := allPartner{}
	err := allPartner.InitAllPartnerInfo()
	if err != nil {
		panic("Init all partner error: " + err.Error())
	}

	return &allPartner
}

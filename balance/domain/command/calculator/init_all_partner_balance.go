package calculator

import (
	"errors"
	"fmt"
	"high-performance-payment-gateway/balance-service/balance/domain/command/logs_request_balance"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"sync"
)

/**
all information partner for calculator balance
*/

type (
	AllPartner struct {
		allPartner        map[string]partnerBalance
		muLock            sync.Mutex
		cnRechargeLog     sql.Connect
		cnBalance         sql.Connect
		logRequestBalance logs_request_balance.LogsInterface
	}

	AllPartnerInterface interface {
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

func (allP *AllPartner) LoadAllPartnerInfo() (map[string]partnerBalance, error) {
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
		cnRechargeLog:         allP.cnRechargeLog,
		cnBalance:             allP.cnBalance,
		logRequestBalance:     allP.logRequestBalance,
	}

	return fake, nil
}

func (allP *AllPartner) InitAllPartnerInfo() error {
	allPartner, err := allP.LoadAllPartnerInfo()
	if err != nil {
		return err
	}

	balancePlaceHolderHistory := NewBalancePlaceHolderHistory()

	allP.muLock.Lock()
	//init raw AllPartner
	for k, v := range allPartner {
		allP.allPartner[k] = v
	}

	//merger balancePlaceHolderHistory to AllPartner
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

func (allP *AllPartner) UpdateOnePartner(p partnerBalance) error {
	p.muLock.Lock()
	key := allP.getKeyOnePartner(p)
	allP.allPartner[key] = p
	return nil
}

func (allP *AllPartner) getKeyOnePartner(p partnerBalance) string {
	return p.partnerCode
}

func (allP *AllPartner) GetOnePartner(partnerCode string) (partnerBalance, error) {
	partner, ok := allP.allPartner[partnerCode]
	if !ok {
		err := fmt.Sprintf("partnercode %s not exists", partnerCode)
		return partnerBalance{}, errors.New(err)
	}

	return partner, nil
}

func InitAllPartnerData() AllPartnerInterface {
	allPartner := AllPartner{}
	err := allPartner.InitAllPartnerInfo()
	if err != nil {
		panic("Init all partner error: " + err.Error())
	}

	return &allPartner
}

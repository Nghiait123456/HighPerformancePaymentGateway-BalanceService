package main

import (
	"errors"
	"fmt"
	"sync"
)

/**
all information partner for calculator balance
*/

var EmergencyS = NewEmergencyStop()

func EmergencyStop() emergencyStopInterface {
	return EmergencyS
}

type allPartner struct {
	allPartner map[string]partnerBalance
	muLock     sync.Mutex
}

type allPartnerInterface interface {
	initPartnersInterface
}

type initPartnersInterface interface {
	LoadAllPartnerInfo() (map[string]partnerBalance, error)
	InitAllPartnerInfo() error
	UpdateOnePartner(p partnerBalance) error
	GetOnePartner(partnerCode string) (partnerBalance, error)
	getKeyOnePartner(p partnerBalance) string
}

func (allP *allPartner) LoadAllPartnerInfo() (map[string]partnerBalance, error) {
	fake := make(map[string]partnerBalance)
	fake["partner_test"] = partnerBalance{
		partnerCode:           "partner_test",
		partnerName:           "TEST",
		partnerIdentification: 1,
		amountTotal:           1000000,
		amountPlaceHolder:     4,
		status:                "active",
		muLock:                sync.Mutex{},
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

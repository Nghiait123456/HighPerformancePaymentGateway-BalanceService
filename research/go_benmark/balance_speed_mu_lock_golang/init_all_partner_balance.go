package main

import (
	"errors"
	"fmt"
	"sync"
)

/**
all information partner for calculator balance
*/

type allPartner struct {
	allPartner map[string]partnerInfo
	muLock     sync.Mutex
}

type allPartnerInterface interface {
	initPartnersInterface
}

type initPartnersInterface interface {
	LoadAllPartnerInfo() (map[string]partnerInfo, error)
	InitAllPartnerInfo() error
	UpdateOnePartner(p partnerInfo) error
	GetOnePartner(partnerCode string) (partnerInfo, error)
	getKeyOnePartner(p partnerInfo) string
}

func (allP *allPartner) LoadAllPartnerInfo() (map[string]partnerInfo, error) {
	fake := make(map[string]partnerInfo)
	fake["partner_test"] = partnerInfo{
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

func (allP *allPartner) UpdateOnePartner(p partnerInfo) error {
	p.muLock.Lock()
	key := allP.getKeyOnePartner(p)
	allP.allPartner[key] = p
	return nil
}

func (allP *allPartner) getKeyOnePartner(p partnerInfo) string {
	return p.partnerCode
}

func (allP *allPartner) GetOnePartner(partnerCode string) (partnerInfo, error) {
	partner, ok := allP.allPartner[partnerCode]
	if !ok {
		err := fmt.Sprintf("partnercode %s not exists", partnerCode)
		return partnerInfo{}, errors.New(err)
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

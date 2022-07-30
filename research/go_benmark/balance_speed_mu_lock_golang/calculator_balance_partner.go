package main

import (
	"errors"
	"fmt"
	"sync"
)

type balancerRequest struct {
	amountRequest         uint
	partnerCode           string
	partnerIdentification uint
	// create order, update amount when partner recharge
	typeRequest string
}

var (
	typeRequestPayment  = "payment"
	typeRequestRecharge = "recharge"
)

type partnerInfo struct {
	partnerCode           string
	partnerName           string
	partnerIdentification uint
	amountTotal           uint
	amountPlaceHolder     uint
	status                string
	muLock                sync.Mutex
}

type calculatorBalancer interface {
	isValidAmount() bool
	isAccessOrder(amountRequest uint, typeRequest string) bool
	increaseAmount(amountRequest uint) error
	increaseAmountPlaceHolder(amountRequest uint) error
	decreaseAmount(amountRequest uint) error
	decreaseAmountPlaceHolder(amountRequest uint) error
}

func (pI *partnerInfo) isValidAmount() bool {
	return pI.amountTotal >= 0 && pI.amountPlaceHolder >= 0 && pI.amountTotal >= pI.amountPlaceHolder
}

func (pI *partnerInfo) isAccessOrder(amountRequest uint, typeRequest string) bool {
	switch typeRequest {
	case typeRequest:
		if amountRequest+pI.amountPlaceHolder <= pI.amountTotal {
			return true
		}
		return false

	case typeRequestRecharge:
		return true

	default:
		err := fmt.Sprintf("typeRequest %s to balancer service not exits", typeRequest)
		panic(err)
	}
}

func (pI *partnerInfo) increaseAmount(amountRequest uint) error {
	pI.amountTotal += amountRequest
	return nil
}

func (pI *partnerInfo) increaseAmountPlaceHolder(amountRequest uint) error {
	pI.amountPlaceHolder += amountRequest
	return nil
}

func (pI *partnerInfo) decreaseAmount(amountRequest uint) error {
	if pI.amountTotal < amountRequest {
		err := fmt.Sprintf("amountRequest %i greater than amountTatal %i in partnerCode %s", amountRequest, pI.amountTotal, pI.partnerCode)
		return errors.New(err)
	}

	pI.amountTotal -= amountRequest
	return nil
}

func (pI *partnerInfo) decreaseAmountPlaceHolder(amountRequest uint) error {
	pI.amountPlaceHolder -= amountRequest
	return nil
}

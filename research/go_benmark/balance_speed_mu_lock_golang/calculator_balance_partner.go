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

type partnerBalance struct {
	partnerCode           string
	partnerName           string
	partnerIdentification uint
	amountTotal           uint
	amountPlaceHolder     uint
	status                string
	muLock                sync.Mutex
}

type calculatorBalancerInterface interface {
	isValidAmount() bool
	isApproveOrder(b balancerRequest) (bool, string)
	increaseAmount(amountRequest uint) error
	increaseAmountPlaceHolder(amountRequest uint) error
	decreaseAmount(amountRequest uint) error
	decreaseAmountPlaceHolder(amountRequest uint) error
	HandleOneRequestBalance(b balancerRequest) (bool, error)
	updateRequestApprovedLocalInMemory(b balancerRequest)
	saveLogsPlaceHolder(b balancerRequest) (bool, error)
	updateRequestApprovedDB(b balancerRequest) (bool, error)
	saveLogsAmountReCharge(b balancerRequest) (bool, error)
}

func (pB *partnerBalance) isValidAmount() bool {
	return pB.amountTotal >= 0 && pB.amountPlaceHolder >= 0 && pB.amountTotal >= pB.amountPlaceHolder
}

func (pB *partnerBalance) isApproveOrder(b balancerRequest) (bool, string) {
	switch b.typeRequest {
	case typeRequestPayment:
		if b.amountRequest+pB.amountPlaceHolder <= pB.amountTotal {
			return true, ""
		}
		return false, "not enough money"

	case typeRequestRecharge:
		return true, ""

	default:
		err := fmt.Sprintf("typeRequest %s to balancer service not exits", b.typeRequest)
		EmergencyStop().ThrowEmergencyStop()
		panic(err)
	}
}

func (pB *partnerBalance) increaseAmount(amountRequest uint) {
	pB.amountTotal += amountRequest
}

func (pB *partnerBalance) increaseAmountPlaceHolder(amountRequest uint) {
	pB.amountPlaceHolder += amountRequest
}

func (pB *partnerBalance) decreaseAmount(amountRequest uint) error {
	if pB.amountTotal < amountRequest {
		err := fmt.Sprintf("amountRequest %i greater than amountTatal %i in partnerCode %s", amountRequest, pB.amountTotal, pB.partnerCode)
		return errors.New(err)
	}

	pB.amountTotal -= amountRequest
	return nil
}

func (pB *partnerBalance) decreaseAmountPlaceHolder(amountRequest uint) error {
	if pB.amountPlaceHolder < amountRequest {
		err := fmt.Sprintf("amountRequest %i greater than amountPlaceHolder %i in partnerCode %s", amountRequest, pB.amountPlaceHolder, pB.partnerCode)
		return errors.New(err)
	}

	pB.amountPlaceHolder -= amountRequest
	return nil
}

// HandleOneRequestBalance is endpoint call check all process
func (pB *partnerBalance) HandleOneRequestBalance(b balancerRequest) (bool, error) {
	pB.muLock.Lock()
	if EmergencyStop().IsStop() {
		return true, nil
	}

	if !pB.isValidAmount() {
		pB.muLock.Unlock()
		return false, errors.New("amount partner not valid")
	}

	approved, errA := pB.isApproveOrder(b)
	if !approved {
		pB.muLock.Unlock()
		return false, errors.New(errA)
	}

	//update to local in memory
	pB.updateRequestApprovedLocalInMemory(b)

	updatedDB, errUDB := pB.updateRequestApprovedDB(b)
	if !updatedDB {
		return false, errUDB
	}

	//release lock for performance
	pB.muLock.Unlock()

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestPaymentLocalInMemory(b balancerRequest) {
	// update local in memory
	pB.increaseAmountPlaceHolder(b.amountRequest)
}

func (pB *partnerBalance) updateTypeRequestPaymentDB(b balancerRequest) (bool, error) {
	// save log place holder
	saveLog1, errSL1 := pB.saveLogsPlaceHolder(b)
	if !saveLog1 {
		//try again
		saveLog2, errSL2 := pB.saveLogsPlaceHolder(b)
		if !saveLog2 {
			//EmergencyStop().ThrowEmergencyStop()
			err := fmt.Sprintf("updateTypeRequestPaymentDB error, err1: %s, err2 : %s ", errSL1.Error(), errSL2.Error())
			panic(err)
		}
	}

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestRechargeLocalInMemory(b balancerRequest) {
	// update local in memory
	pB.increaseAmount(b.amountRequest)
}

func (pB *partnerBalance) updateTypeRequestRechargeDB(b balancerRequest) (bool, error) {
	// save log place holder
	saveLog1, errSL1 := pB.saveLogsAmountReCharge(b)
	if !saveLog1 {
		//try again
		saveLog2, errSL2 := pB.saveLogsAmountReCharge(b)
		if !saveLog2 {
			EmergencyStop().ThrowEmergencyStop()
			err := fmt.Sprintf("updateTypeRequestRechargeDB error, err1: %s, err2: %s", errSL1.Error(), errSL2.Error())
			panic(err)
		}
	}

	return true, nil
}

func (pB *partnerBalance) updateRequestApprovedLocalInMemory(b balancerRequest) {
	switch b.typeRequest {
	case typeRequestPayment:
		pB.updateTypeRequestPaymentLocalInMemory(b)
	case typeRequestRecharge:
		pB.updateTypeRequestRechargeLocalInMemory(b)

	default:
		err := fmt.Sprintf("typeRequest %s to balancer service not exits", b.typeRequest)
		panic(err)
	}
}

func (pB *partnerBalance) updateRequestApprovedDB(b balancerRequest) (bool, error) {
	switch b.typeRequest {
	case typeRequestPayment:
		return pB.updateTypeRequestPaymentDB(b)
	case typeRequestRecharge:
		return pB.updateTypeRequestRechargeDB(b)

	default:
		err := fmt.Sprintf("typeRequest %s to balancer service not exits", b.typeRequest)
		panic(err)
	}
}

func (pB *partnerBalance) saveLogsPlaceHolder(b balancerRequest) (bool, error) {
	// todo save logs to DB logs requet balance, begin transaction, db is sharding
	return true, nil
}

// SaveLogsAmountReCharge save in same DB with totalAmount Balance
func (pB *partnerBalance) saveLogsAmountReCharge(b balancerRequest) (bool, error) {
	// todo start transacions, update amount and update logs
	return true, nil
}

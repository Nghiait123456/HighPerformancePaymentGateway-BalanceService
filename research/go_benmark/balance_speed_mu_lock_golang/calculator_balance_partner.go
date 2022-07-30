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
	updateRequestApproved(b balancerRequest) (bool, error)
	saveLogsPlaceHolder(b balancerRequest) (bool, error)
	updateTypeRequestRecharge(b balancerRequest) (bool, error)
	SaveLogsAmountReCharge(b balancerRequest) (bool, error)
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

func (pB *partnerBalance) increaseAmount(amountRequest uint) error {
	pB.amountTotal += amountRequest
	return nil
}

func (pB *partnerBalance) increaseAmountPlaceHolder(amountRequest uint) error {
	pB.amountPlaceHolder += amountRequest

	return nil
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
	updatedInMemory, errUI := pB.updateRequestApprovedLocalInMemory(b)
	if !updatedInMemory {
		return false, errUI
	}
	//release lock for performance
	pB.muLock.Unlock()

	updatedDB, errUDB := pB.updateRequestApprovedDB(b)
	if !updatedDB {
		return false, errUDB
	}

}

func (pB *partnerBalance) updateTypeRequestPaymentLocalInMemory(b balancerRequest) (bool, error) {
	// update local in memory
	err := pB.increaseAmountPlaceHolder(b.amountRequest)
	if err != nil {
		//rollback local in memory
		err := pB.decreaseAmountPlaceHolder(b.amountRequest)
		if err != nil {
			panic("don't roollback amount")
		}
		return false, err
	}

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestPaymentDB(b balancerRequest) (bool, error) {
	// save log place holder
	saveLog, errSL := pB.saveLogsPlaceHolder(b)
	if !saveLog {
		//rollback  local in memory
		pB.muLock.Lock()
		err := pB.decreaseAmountPlaceHolder(b.amountRequest)
		if err != nil {
			pB.muLock.Unlock()
			EmergencyStop().ThrowEmergencyStop()
			panic("don't roollback amount place holder")
		}

		pB.muLock.Unlock()
		return false, errSL
	}

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestRechargeLocalInMemory(b balancerRequest) (bool, error) {
	// update local in memory
	err := pB.increaseAmount(b.amountRequest)
	if err != nil {
		//rollback local in memory
		err := pB.decreaseAmount(b.amountRequest)
		if err != nil {
			panic("don't roollback amount")
		}

		return false, err
	}

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestRechargeDB(b balancerRequest) (bool, error) {
	// save log place holder
	saveLog, errSL := pB.saveLogsAmountReCharge(b)
	if !saveLog {
		//rollback  local in memory
		pB.muLock.Lock()
		err := pB.decreaseAmount(b.amountRequest)
		if err != nil {
			pB.muLock.Unlock()
			EmergencyStop().ThrowEmergencyStop()
			panic("don't roollback amount")
		}

		pB.muLock.Unlock()
		return false, errSL
	}

	return true, nil
}

func (pB *partnerBalance) updateRequestApprovedLocalInMemory(b balancerRequest) (bool, error) {
	switch b.typeRequest {
	case typeRequestPayment:
		return pB.updateTypeRequestPaymentLocalInMemory(b)
	case typeRequestRecharge:
		return pB.updateTypeRequestRechargeLocalInMemory(b)

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
	// todo save logs to DB logs requet balance, db is sharding
	return true, nil
}

// SaveLogsAmountReCharge save in same DB with totalAmount Balance
func (pB *partnerBalance) saveLogsAmountReCharge(b balancerRequest) (bool, error) {
	// todo start transacions, update amount and update logs
	return true, nil
}

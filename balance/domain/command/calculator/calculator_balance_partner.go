package calculator

import (
	"errors"
	"fmt"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/shard_balance_logs"
	"sync"
	"time"
)

type (
	balancerRequest struct {
		amountRequest         uint64
		partnerCode           string
		partnerIdentification uint
		orderID               uint64
		// create order, update amount when partner recharge
		typeRequest string
	}

	partnerBalance struct {
		partnerCode           string
		partnerName           string
		partnerIdentification uint
		amountTotal           uint64
		amountPlaceHolder     uint64
		status                string
		muLock                sync.Mutex
		EStop                 emergencyStopInterface
		lbShardBalanceLog     shard_balance_logs.LBShardLogInterface
	}

	saveLogsDB struct {
		b                 balancerRequest
		pb                partnerBalance
		lbShardBalanceLog shard_balance_logs.LBShardLogInterface
		saveLog           shard_balance_logs.SaveLogInterface
	}
	partnerBalanceInterface interface {
		calculatorPartnerBalancerInterface
	}

	calculatorPartnerBalancerInterface interface {
		isValidAmount() bool
		isApproveOrder(b balancerRequest) (bool, string)
		increaseAmount(amountRequest uint64)
		increaseAmountPlaceHolder(amountRequest uint64)
		decreaseAmount(amountRequest uint64) error
		decreaseAmountPlaceHolder(amountRequest uint64) error
		HandleOneRequestBalance(b balancerRequest) (bool, error)
		updateRequestApprovedLocalInMemory(b balancerRequest)
		saveLogsPlaceHolder(u saveLogsDB) (bool, error)
		saveRequestApprovedDB(u saveLogsDB) (bool, error)
		saveLogsAmountReCharge(u saveLogsDB) (bool, error)
	}
)

var (
	typeRequestPayment  = "payment"
	typeRequestRecharge = "recharge"
)

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
		pB.EStop.ThrowEmergencyStop()
		panic(err)
	}
}

func (pB *partnerBalance) increaseAmount(amountRequest uint64) {
	pB.amountTotal += amountRequest
}

func (pB *partnerBalance) increaseAmountPlaceHolder(amountRequest uint64) {
	pB.amountPlaceHolder += amountRequest
}

func (pB *partnerBalance) decreaseAmount(amountRequest uint64) error {
	if pB.amountTotal < amountRequest {
		err := fmt.Sprintf("amountRequest %i greater than amountTatal %i in partnerCode %s", amountRequest, pB.amountTotal, pB.partnerCode)
		return errors.New(err)
	}

	pB.amountTotal -= amountRequest
	return nil
}

func (pB *partnerBalance) decreaseAmountPlaceHolder(amountRequest uint64) error {
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
	if pB.EStop.IsStop() {
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

	saveLogsDB := saveLogsDB{
		b:                 b,
		pb:                pB.ValueObject(),
		saveLog:           shard_balance_logs.SaveLog{},
		lbShardBalanceLog: pB.lbShardBalanceLog,
	}
	
	updatedDB, errUDB := pB.saveRequestApprovedDB(saveLogsDB)
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

func (pB *partnerBalance) ValueObject() partnerBalance {
	pBNews := *pB
	pBNews.muLock = sync.Mutex{}
	pBNews.EStop = nil
	return pBNews
}

func (pB *partnerBalance) saveTypeRequestPaymentDB(u saveLogsDB) (bool, error) {
	// save log place holder
	status, err := pB.saveLogsPlaceHolder(u)
	return status, err
}

func (pB *partnerBalance) updateTypeRequestRechargeLocalInMemory(b balancerRequest) {
	// update local in memory
	pB.increaseAmount(b.amountRequest)
}

func (pB *partnerBalance) saveTypeRequestRechargeDB(s saveLogsDB) (bool, error) {
	// save log place holder
	status, err := pB.saveLogsAmountReCharge(s)
	return status, err
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

func (pB *partnerBalance) saveRequestApprovedDB(s saveLogsDB) (bool, error) {
	switch s.b.typeRequest {
	case typeRequestPayment:
		return pB.saveTypeRequestPaymentDB(s)
	case typeRequestRecharge:
		return pB.saveTypeRequestRechargeDB(s)

	default:
		err := fmt.Sprintf("typeRequest %s to balancer service not exits", s.b.typeRequest)
		panic(err)
	}
}

func (pB *partnerBalance) saveLogsPlaceHolder(s saveLogsDB) (bool, error) {
	var bl orm.BalanceLog

	bl.OrderId = s.b.orderID
	bl.PartnerCode = s.b.partnerCode
	bl.AmountRequest = s.b.amountRequest
	bl.AmountPlaceHolder = s.pb.amountPlaceHolder
	bl.Balance = s.pb.amountTotal
	bl.Status = orm.BalanceLog{}.StatusProcessing()
	bl.CreatedAt = uint32(time.Now().Unix())
	bl.UpdatedAt = uint32(time.Now().Unix())

	ok := s.saveLog.Save(s.lbShardBalanceLog, bl)
	if ok != nil {
		return false, ok
	}

	return true, nil
}

// SaveLogsAmountReCharge save in same DB with totalAmount Balance
func (pB *partnerBalance) saveLogsAmountReCharge(u saveLogsDB) (bool, error) {
	// todo start transacions, update amount and update logs
	return true, nil
}

func NewPartnerBalance() partnerBalanceInterface {
	return &partnerBalance{}
}

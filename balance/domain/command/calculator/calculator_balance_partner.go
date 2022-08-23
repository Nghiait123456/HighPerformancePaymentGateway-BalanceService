package calculator

import (
	"errors"
	"fmt"
	"high-performance-payment-gateway/balance-service/balance/domain/command/logs_request_balance"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
	"os"
	"sync"
	"time"
)

type (
	BalancerRequest struct {
		AmountRequest         uint64
		PartnerCode           string
		partnerIdentification uint
		orderID               uint64
		// create order, update amount when partner recharge
		typeRequest string
	}

	partnerBalance struct {
		partnerCode           string
		partnerName           string
		partnerIdentification uint
		balance               uint64
		amountPlaceHolder     uint64
		indexLogRequestLatest uint64
		status                string
		muLock                sync.Mutex
		EStop                 emergencyStop
		//lbShardBalanceLog     shard_balance_logs.LBShardLogInterface
		cnRechargeLog     sql.Connect
		cnBalance         sql.Connect
		logRequestBalance logs_request_balance.LogsInterface
	}

	saveLogsDB struct {
		b  BalancerRequest
		pb partnerBalance
		//lbShardBalanceLog     shard_balance_logs.LBShardLogInterface
		//saveLogRequestBalance shard_balance_logs.SaveLogRequestBalanceInterface
	}
	partnerBalanceInterface interface {
		calculatorPartnerBalancerInterface
	}

	calculatorPartnerBalancerInterface interface {
		isValidAmount() bool
		isFull() bool
		isApproveOrder(b BalancerRequest) (bool, string)
		increaseAmount(amountRequest uint64)
		increaseAmountPlaceHolder(amountRequest uint64)
		decreaseAmount(amountRequest uint64) error
		decreaseAmountPlaceHolder(amountRequest uint64) error
		increaseIndexLogRequestLatest()
		HandleOneRequestBalance(b BalancerRequest) (bool, error)
		commitPlaceHolderToLocalInMemory() error
		rollbackCommitPlaceHolderToLocalInMemory() error
		updateRequestApprovedLocalInMemory(b BalancerRequest)
		revertRequestApprovedLocalInMemory(b BalancerRequest) error
		saveLogsPlaceHolder(s saveLogsDB) (bool, error)
		saveRequestApprovedDB(s saveLogsDB) (bool, error)
		saveLogsAndAmountReChargeDB(s saveLogsDB) (bool, error)
	}
)

var (
	typeRequestPayment  = "payment"
	typeRequestRecharge = "recharge"
)

func (pB *partnerBalance) isFull() bool {
	return pB.balance == pB.amountPlaceHolder
}

func (pB *partnerBalance) isValidAmount() bool {
	return pB.balance >= 0 && pB.amountPlaceHolder >= 0 && pB.balance > pB.amountPlaceHolder
}

func (pB *partnerBalance) isApproveOrder(b BalancerRequest) (bool, string) {
	switch b.typeRequest {
	case typeRequestPayment:
		if b.AmountRequest+pB.amountPlaceHolder <= pB.balance {
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
	pB.balance += amountRequest
}

func (pB *partnerBalance) increaseAmountPlaceHolder(amountRequest uint64) {
	pB.amountPlaceHolder += amountRequest
}

func (pB *partnerBalance) increaseIndexLogRequestLatest() {
	pB.indexLogRequestLatest += 1
}

func (pB *partnerBalance) decreaseIndexLogRequestLatest() error {
	if pB.indexLogRequestLatest < 1 {
		err := fmt.Sprintf("indexLogRequestLatest %d greater than zero when decrease", pB.indexLogRequestLatest)
		return errors.New(err)
	}

	pB.indexLogRequestLatest -= 1
	return nil
}

func (pB *partnerBalance) decreaseAmount(amountRequest uint64) error {
	if pB.balance < amountRequest {
		err := fmt.Sprintf("AmountRequest %i greater than amountTatal %i in PartnerCode %s", amountRequest, pB.balance, pB.partnerCode)
		return errors.New(err)
	}

	pB.balance -= amountRequest
	return nil
}

func (pB *partnerBalance) decreaseAmountPlaceHolder(amountRequest uint64) error {
	if pB.amountPlaceHolder < amountRequest {
		err := fmt.Sprintf("AmountRequest %i greater than amountPlaceHolder %i in PartnerCode %s", amountRequest, pB.amountPlaceHolder, pB.partnerCode)
		return errors.New(err)
	}

	pB.amountPlaceHolder -= amountRequest
	return nil
}

// HandleOneRequestBalance is endpoint call check all process
func (pB *partnerBalance) HandleOneRequestBalance(b BalancerRequest) (bool, error) {
	pB.muLock.Lock()
	defer pB.muLock.Unlock()

	if pB.EStop.IsStop() {
		return false, errors.New("partner is stopping")
	}

	if pB.isFull() {
		return false, errors.New("balance is full, please try again")
	}

	if !pB.isValidAmount() {
		return false, errors.New("amount partner not valid")
	}

	approved, errA := pB.isApproveOrder(b)
	if !approved {
		return false, errors.New(errA)
	}

	//update to local in memory
	pB.updateRequestApprovedLocalInMemory(b)

	saveLogsDB := saveLogsDB{
		b:  b,
		pb: pB.ValueObject(),
	}

	updatedDB, errUDB := pB.saveRequestApprovedDB(saveLogsDB)
	if !updatedDB {
		//revert if error
		errRv := pB.revertRequestApprovedLocalInMemory(b)
		if errRv != nil {
			panic(fmt.Sprintf("don't revert request approved ib local in memory, err: %s", errRv.Error()))
			os.Exit(0)
		}
		return false, errUDB
	}

	return true, nil
}

func (pB *partnerBalance) updateTypeRequestPaymentLocalInMemory(b BalancerRequest) {
	// update local in memory
	pB.increaseAmountPlaceHolder(b.AmountRequest)
	pB.increaseIndexLogRequestLatest()
}

func (pB *partnerBalance) commitPlaceHolderToLocalInMemory() error {
	if pB.amountPlaceHolder <= 0 {
		errM := fmt.Sprintf("balanceIncrement must is unsigned")
		panic(errM)
		os.Exit(0)
	}

	if pB.balance < pB.amountPlaceHolder {
		errM := fmt.Sprintf("balance %d greater amountPlaceHolder %d with PartnerCode %s", pB.balance, pB.amountPlaceHolder, pB.partnerCode)
		panic(errM)
		os.Exit(0)
	}

	//update balance
	pB.balance -= pB.amountPlaceHolder
	return nil
}

func (pB *partnerBalance) rollbackCommitPlaceHolderToLocalInMemory() error {
	if pB.amountPlaceHolder <= 0 {
		errM := fmt.Sprintf("balanceIncrement must is unsigned")
		panic(errM)
		os.Exit(0)
	}

	if pB.balance < pB.amountPlaceHolder {
		errM := fmt.Sprintf("balance %d greater amountPlaceHolder %d with PartnerCode %s", pB.balance, pB.amountPlaceHolder, pB.partnerCode)
		panic(errM)
		os.Exit(0)
	}

	//update balance
	pB.balance += pB.amountPlaceHolder
	return nil
}

func (pB *partnerBalance) revertTypeRequestPaymentLocalInMemory(b BalancerRequest) error {
	// update balance
	errAP := pB.decreaseAmountPlaceHolder(b.AmountRequest)
	if errAP != nil {
		return errAP
	}

	//update latest index log
	errILR := pB.decreaseIndexLogRequestLatest()
	if errILR != nil {
		return errILR
	}

	return nil
}

func (pB *partnerBalance) ValueObject() partnerBalance {
	pBValue := partnerBalance{}
	pBValue.partnerCode = pB.partnerCode
	pBValue.partnerName = pB.partnerName
	pBValue.partnerIdentification = pB.partnerIdentification
	pBValue.amountPlaceHolder = pB.amountPlaceHolder
	pBValue.status = pB.status

	return pBValue
}

func (pB *partnerBalance) saveTypeRequestPaymentDB(s saveLogsDB) (bool, error) {
	// save log place holder
	status, err := pB.saveLogsPlaceHolder(s)
	return status, err
}

func (pB *partnerBalance) updateTypeRequestRechargeLocalInMemory(b BalancerRequest) {
	// update local in memory
	pB.increaseAmount(b.AmountRequest)
}

func (pB *partnerBalance) revertTypeRequestRechargeLocalInMemory(b BalancerRequest) error {
	// update local in memory
	return pB.decreaseAmount(b.AmountRequest)
}

func (pB *partnerBalance) saveTypeRequestRechargeDB(s saveLogsDB) (bool, error) {
	// save log recharge
	status, err := pB.saveLogsAndAmountReChargeDB(s)
	return status, err
}

func (pB *partnerBalance) updateRequestApprovedLocalInMemory(b BalancerRequest) {
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

func (pB *partnerBalance) revertRequestApprovedLocalInMemory(b BalancerRequest) error {
	switch b.typeRequest {
	case typeRequestPayment:
		return pB.revertTypeRequestPaymentLocalInMemory(b)
	case typeRequestRecharge:
		return pB.revertTypeRequestRechargeLocalInMemory(b)

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
	//var brl orm.BalanceRequestLog
	//brl.OrderId = s.b.orderID
	//brl.PartnerCode = s.b.PartnerCode
	//brl.AmountRequest = s.b.AmountRequest
	//brl.AmountPlaceHolder = s.pb.amountPlaceHolder
	//brl.Balance = s.pb.balance
	//brl.Status = orm.BalanceRequestLog{}.StatusProcessing()
	//brl.CreatedAt = uint32(time.Now().Unix())
	//brl.UpdatedAt = uint32(time.Now().Unix())
	//
	//ok := s.saveLogRequestBalance.Save(s.lbShardBalanceLog, brl)
	//if ok != nil {
	//	return false, ok
	//}

	o := logs_request_balance.OrderLog{
		OrderId: s.b.orderID,
		Amount:  s.b.AmountRequest,
		Status:  "processing",
	}

	pB.logRequestBalance.Push(o)
	return true, nil
}

// SaveLogsAmountReCharge save in same DB with totalAmount Balance
func (pB *partnerBalance) saveLogsAndAmountReChargeDB(s saveLogsDB) (bool, error) {
	var brl orm.BalanceRechargeLog
	var u repository.UpdateLogsAndBalance
	rp := repository.NewBalanceAndBalanceRequestLogRepository()

	// BalanceRechargeLog
	brl.OrderId = s.b.orderID
	brl.PartnerCode = s.b.PartnerCode
	brl.AmountRecharge = s.b.AmountRequest
	brl.Balance = s.pb.balance
	brl.Status = orm.BalanceRechargeLog{}.StatusSuccess()
	brl.CreatedAt = uint32(time.Now().Unix())
	brl.UpdatedAt = uint32(time.Now().Unix())

	// init update data
	u.PartnerCode = s.pb.partnerCode
	u.Balance = s.pb.balance
	u.Log = brl
	u.CN = pB.cnRechargeLog

	err := rp.SaveLogsAndUpdateBalanceReChargeDB(u)
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewPartnerBalance() partnerBalanceInterface {
	return &partnerBalance{}
}

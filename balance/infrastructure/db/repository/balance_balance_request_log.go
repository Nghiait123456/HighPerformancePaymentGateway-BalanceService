package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceAndBalanceRequestLog struct {
		BaseRepo BaseInterface
	}

	UpdateLogsAndBalance struct {
		PartnerCode string
		Balance     uint64
		Log         orm.BalanceRechargeLog
		CN          sql.Connect
	}

	// implement query in trans with 2 table orm  orm.Balance and orm.BalanceRechargeLog
	BalanceAndBalanceRequestLogInterface interface {
		SetConnect(cn sql.Connect)
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		GetConnectFrGlobalCf() sql.Connect
		ResetContext()
		SaveLogsAndUpdateBalanceReChargeDB(u UpdateLogsAndBalance) error
	}
)

func (rp *BalanceAndBalanceRequestLog) SaveLogsAndUpdateBalanceReChargeDB(u UpdateLogsAndBalance) error {
	rpBRL := NewBalanceRechargeLogRepository()
	rpBL := NewBalanceRepository()

	rpBRL.SetConnect(u.CN)
	rpBL.SetConnect(u.CN)

	rpBL.DB().Begin()
	//update Balance
	errBL := rpBL.UpdateBalanceByPartnerCode(u.PartnerCode, u.Balance)
	if errBL != nil {
		rpBL.DB().Rollback()
		return errBL
	}

	// save Log
	erPBRL := rpBRL.CreateNew(u.Log)
	if erPBRL != nil {
		rpBL.DB().Rollback()
		return erPBRL
	}

	rpBL.DB().Commit()
	return nil
}

func (rp *BalanceAndBalanceRequestLog) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp BalanceAndBalanceRequestLog) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
}

func (rp BalanceAndBalanceRequestLog) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

//SetTimeout ms
func (rp *BalanceAndBalanceRequestLog) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *BalanceAndBalanceRequestLog) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *BalanceAndBalanceRequestLog) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *BalanceAndBalanceRequestLog) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceAndBalanceRequestLogRepository() BalanceAndBalanceRequestLogInterface {
	b := BalanceAndBalanceRequestLog{
		BaseRepo: NewBaseRepository(),
	}
	b.SetConnect(b.GetConnectFrGlobalCf())

	return &b
}

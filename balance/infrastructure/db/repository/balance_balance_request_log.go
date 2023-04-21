package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceAndBalanceRequestLog struct {
		BalanceRequestLogeOrm orm.BalanceRequestLog
		BaseRepo              BaseInterface
	}

	UpdateLogsAndBalance struct {
		PartnerCode string
		Balance     uint64
		Log         orm.BalanceRechargeLog
		CN          sql.Connect
	}

	// implement query in trans with 2 table orm  orm.Balance and orm.BalanceRechargeLog
	BalanceAndBalanceRequestLogInterface interface {
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		SaveLogsAndUpdateBalanceReChargeDB(u UpdateLogsAndBalance) error
	}
)

func (rp *BalanceAndBalanceRequestLog) SaveLogsAndUpdateBalanceReChargeDB(u UpdateLogsAndBalance) error {
	var cn sql.Connect //todo get cn from global config
	rpBRL := NewBalanceRechargeLogRepository(cn)
	rpBL := NewBalanceRepository(u.CN)

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

func NewBalanceAndBalanceRequestLogRepository(cn sql.Connect) BalanceAndBalanceRequestLogInterface {
	baseRepo := NewBaseRepository(cn)
	rp := BalanceAndBalanceRequestLog{
		BaseRepo: baseRepo,
	}

	return &rp
}

package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceRechargeLog struct {
		BalanceRQLogOrm orm.BalanceRechargeLog
		BaseRepo        BaseInterface
	}

	BalanceRechargeLogInterface interface {
		SetConnect(cn sql.Connect)
		GetConnectFrGlobalCf() sql.Connect
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		GetById(id uint32) (orm.BalanceRechargeLog, error)
		CreateNew(brl orm.BalanceRechargeLog) error
		CreateNewWithTrans(brl orm.BalanceRechargeLog) error
		UpdateAllField(update orm.BalanceRechargeLog) error
		UpdateByOrderId(orderId uint64, update map[string]interface{}) error
		UpdateById(id uint32, update map[string]interface{}) error
	}
)

func (rp *BalanceRechargeLog) GetById(id uint32) (orm.BalanceRechargeLog, error) {
	var balanceLog orm.BalanceRechargeLog

	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("id = ?", id).First(&balanceLog)
	if rs.Error != nil {
		return balanceLog, rs.Error
	}

	return balanceLog, nil
}

func (rp *BalanceRechargeLog) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp BalanceRechargeLog) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
}

func (rp BalanceRechargeLog) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

func (rp *BalanceRechargeLog) CreateNew(brl orm.BalanceRechargeLog) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Create(&brl)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceRechargeLog) CreateNewWithTrans(brl orm.BalanceRechargeLog) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	tx := rp.DB().Begin()
	rs := tx.Create(&brl)
	if rs.Error != nil {
		tx.Rollback()
		return rs.Error
	}

	tx.Commit()
	return nil
}

func (rp *BalanceRechargeLog) UpdateAllField(update orm.BalanceRechargeLog) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceRechargeLog) UpdateById(id uint32, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.BalanceRequestLog{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceRechargeLog) UpdateByOrderId(orderId uint64, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.BalanceRechargeLog{}).Where("order_id = ?", orderId).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

//SetTimeout ms
func (rp *BalanceRechargeLog) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *BalanceRechargeLog) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *BalanceRechargeLog) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *BalanceRechargeLog) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceRechargeLogRepository() BalanceRechargeLogInterface {
	rp := BalanceRechargeLog{}
	rp.SetConnect(rp.GetConnectFrGlobalCf())
	return &rp
}

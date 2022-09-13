package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceRequestLog struct {
		BalanceRQLogOrm orm.BalanceRequestLog
		BaseRepo        BaseInterface
	}

	BalanceRequestLogInterface interface {
		SetConnect(cn sql.Connect)
		GetConnectFrGlobalCf() sql.Connect
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		GetById(id uint32) (orm.BalanceRequestLog, error)
		CreateNew(bll orm.BalanceRequestLog) error
		CreateNewWithTrans(bll orm.BalanceRequestLog) error
		UpdateAllField(update orm.BalanceRequestLog) error
		UpdateByOrderId(orderId uint64, update map[string]interface{}) error
		UpdateById(id uint32, update map[string]interface{}) error
	}
)

func (rp *BalanceRequestLog) GetById(id uint32) (orm.BalanceRequestLog, error) {
	var balanceLog orm.BalanceRequestLog

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

func (rp *BalanceRequestLog) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp BalanceRequestLog) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
}

func (rp BalanceRequestLog) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

func (rp *BalanceRequestLog) CreateNew(bll orm.BalanceRequestLog) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceRequestLog) CreateNewWithTrans(bll orm.BalanceRequestLog) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	tx := rp.DB().Begin()
	rs := tx.Create(&bll)
	if rs.Error != nil {
		tx.Rollback()
		return rs.Error
	}

	tx.Commit()
	return nil
}

func (rp *BalanceRequestLog) UpdateAllField(update orm.BalanceRequestLog) error {
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

func (rp *BalanceRequestLog) UpdateById(id uint32, update map[string]interface{}) error {
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

func (rp *BalanceRequestLog) UpdateByOrderId(orderId uint64, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.BalanceRequestLog{}).Where("order_id = ?", orderId).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

//SetTimeout ms
func (rp *BalanceRequestLog) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *BalanceRequestLog) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *BalanceRequestLog) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *BalanceRequestLog) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceLogRepository() BalanceRequestLogInterface {
	rp := BalanceRequestLog{}
	rp.SetConnect(rp.GetConnectFrGlobalCf())
	return &rp
}

package repository

import (
	"context"
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceLog struct {
		DB         *gorm.DB
		BalanceOrm orm.Balance
		BaseRepo   BaseInterface
	}

	BalanceLogInterface interface {
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		GetById(id uint32) (orm.BalanceLog, error)
		CreateNew(bll orm.BalanceLog) error
		UpdateAllField(update orm.BalanceLog) error
		UpdateByOrderId(orderId uint64, update map[string]interface{}) error
		UpdateById(id uint32, update map[string]interface{}) error
	}
)

func (rp *BalanceLog) GetById(id uint32) (orm.BalanceLog, error) {
	var balanceLog orm.BalanceLog

	rp.BaseRepo.UpdateContext(rp.DB)
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB.Where("id = ?", id).First(&balanceLog)
	if rs.Error != nil {
		return balanceLog, rs.Error
	}

	return balanceLog, nil
}

func (rp *BalanceLog) CreateNew(bll orm.BalanceLog) error {
	rp.BaseRepo.UpdateContext(rp.DB)
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB.Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceLog) UpdateAllField(update orm.BalanceLog) error {
	rp.BaseRepo.UpdateContext(rp.DB)
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB.Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceLog) UpdateById(id uint32, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext(rp.DB)
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB.Model(&orm.BalanceLog{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceLog) UpdateByOrderId(orderId uint64, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext(rp.DB)
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB.Model(&orm.BalanceLog{}).Where("order_id = ?", orderId).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

//SetTimeout ms
func (rp *BalanceLog) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *BalanceLog) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *BalanceLog) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *BalanceLog) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceLogRepository() BalanceLogInterface {
	return &BalanceLog{}
}

package repository

import (
	"context"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	Balance struct {
		BalanceOrm orm.Balance
		BaseRepo   BaseInterface
	}

	BalanceInterface interface {
		SetConnect(cn sql.Connect)
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		GetById(id uint32) (orm.Balance, error)
		GetByPartnerCode(partnerCode string) (orm.Balance, error)
		CreateNew(bl orm.Balance) error
		UpdateAllField(update orm.Balance) error
		UpdateByPartnerCode(partnerCode string, update map[string]interface{}) error
		UpdateById(id uint32, update map[string]interface{}) error
	}
)

func (rp *Balance) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp Balance) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

func (rp *Balance) GetById(id uint32) (orm.Balance, error) {
	var balance orm.Balance

	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("id = ?", id).First(&balance)
	if rs.Error != nil {
		return balance, rs.Error
	}

	return balance, nil
}

func (rp *Balance) GetByPartnerCode(partnerCode string) (orm.Balance, error) {
	var balance orm.Balance

	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("partner_code = ?", partnerCode).First(&balance)
	if rs.Error != nil {
		return balance, rs.Error
	}

	return balance, nil
}

func (rp *Balance) CreateNew(bl orm.Balance) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Create(&bl)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *Balance) UpdateAllField(update orm.Balance) error {
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

func (rp *Balance) UpdateByPartnerCode(partnerCode string, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.Balance{}).Where("partner_code = ?", partnerCode).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *Balance) UpdateById(id uint32, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.Balance{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

//SetTimeout ms
func (rp *Balance) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *Balance) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *Balance) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *Balance) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceRepository() BalanceInterface {
	return &Balance{}
}

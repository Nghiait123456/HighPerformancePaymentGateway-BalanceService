package repository

import (
	"context"
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
	"os"
)

type (
	Balance struct {
		BalanceOrm orm.Balance
		BaseRepo   BaseInterface
	}

	CommitAmountPlaceHolder struct {
		PartnerCode           string
		AmountPlaceHolder     uint64
		IndexLogRequestLatest uint64
	}

	BalanceInterface interface {
		SetConnect(cn sql.Connect)
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		GetConnectFrGlobalCf() sql.Connect
		ResetContext()
		GetById(id uint32) (orm.Balance, error)
		GetByPartnerCode(partnerCode string) (orm.Balance, error)
		CreateNew(bl orm.Balance) error
		UpdateAllField(update orm.Balance) error
		UpdateByPartnerCode(partnerCode string, update map[string]interface{}) error
		UpdateBalanceByPartnerCode(partnerCode string, balance uint64) error
		UpdateById(id uint32, update map[string]interface{}) error
		CommitAmountPlaceHolderToBalance(c CommitAmountPlaceHolder) error
	}
)

func (rp *Balance) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp Balance) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
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

func (rp *Balance) CommitAmountPlaceHolderToBalance(c CommitAmountPlaceHolder) error {
	var bl orm.Balance
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	if c.AmountPlaceHolder <= 0 {
		errM := fmt.Sprintf("balanceIncrement must is unsigned")
		panic(errM)
		os.Exit(0)
	}

	rs := rp.DB().Where("partner_code = ?", c.PartnerCode).First(&bl)
	if rs.Error != nil {
		return rs.Error
	}

	if bl.ID == 0 {
		errM := fmt.Sprintf("don't find balance with partnerCode %s", c.PartnerCode)
		panic(errM)
		os.Exit(0)
	}

	if bl.Balance < c.AmountPlaceHolder {
		errM := fmt.Sprintf("balance %d greater amountPlaceHolder %d with partnerCode %s", bl.Balance, c.AmountPlaceHolder, c.PartnerCode)
		panic(errM)
		os.Exit(0)
	}

	// calculator new balance
	bl.Balance -= c.AmountPlaceHolder
	// update index latest
	bl.IndexLogRequestLatest = c.IndexLogRequestLatest

	//commit
	rp.DB().Begin()
	rsU := rp.DB().Updates(&bl)
	if rs.Error != nil {
		rp.DB().Rollback()
		return rsU.Error
	}

	rp.DB().Commit()
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

func (rp *Balance) UpdateBalanceByPartnerCode(partnerCode string, balance uint64) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.Balance{}).Where("partner_code = ?", partnerCode).Update("balance", balance)
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
	b := Balance{}
	b.SetConnect(b.GetConnectFrGlobalCf())
	return &b
}

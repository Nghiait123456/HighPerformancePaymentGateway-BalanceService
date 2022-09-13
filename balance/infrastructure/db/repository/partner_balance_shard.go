package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/entity"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	PartnerBalanceShard struct {
		PartnerBalanceOrm orm.PartnerBalanceShard
		BaseRepo          BaseInterface
	}

	PartnerBalanceShardInterface interface {
		SetConnect(cn sql.Connect)
		GetConnectFrGlobalCf() sql.Connect
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		GetById(id uint32) (orm.PartnerBalanceShard, error)
		CreateNew(bll orm.PartnerBalanceShard) error
		UpdateAllField(update orm.PartnerBalanceShard) error
		UpdateById(id uint32, update map[string]interface{}) error
		GetAllActiveByPartner(partnerCode string) (entity.PartnerBalanceShards, error)
	}
)

func (rp *PartnerBalanceShard) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp PartnerBalanceShard) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
}

func (rp PartnerBalanceShard) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

func (rp *PartnerBalanceShard) GetById(id uint32) (orm.PartnerBalanceShard, error) {
	var shard orm.PartnerBalanceShard

	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("id = ?", id).First(&shard)
	if rs.Error != nil {
		return shard, rs.Error
	}

	return shard, nil
}

func (rp *PartnerBalanceShard) CreateNew(bll orm.PartnerBalanceShard) error {
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

func (rp *PartnerBalanceShard) UpdateAllField(update orm.PartnerBalanceShard) error {
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

func (rp *PartnerBalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.PartnerBalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *PartnerBalanceShard) GetAllActiveByPartner(partnerCode string) (entity.PartnerBalanceShards, error) {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	var rs []orm.PartnerBalanceShard
	r := rp.DB().Where("partner_code = ?", partnerCode).Where("status = ?", rp.PartnerBalanceOrm.StatusActive()).Find(&rs)
	if r.Error != nil {
		return rs, r.Error
	}

	return rs, nil
}

//SetTimeout ms
func (rp *PartnerBalanceShard) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *PartnerBalanceShard) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *PartnerBalanceShard) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *PartnerBalanceShard) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewPartnerBalanceShardRepository() PartnerBalanceShardInterface {
	rp := PartnerBalanceShard{}
	rp.SetConnect(rp.GetConnectFrGlobalCf())
	return &rp
}

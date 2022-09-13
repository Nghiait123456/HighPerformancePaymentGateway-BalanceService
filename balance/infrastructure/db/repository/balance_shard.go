package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type (
	BalanceShard struct {
		BalanceShardOrm orm.BalanceShard
		BaseRepo        BaseInterface
	}

	BalanceShardInterface interface {
		SetConnect(cn sql.Connect)
		GetConnectFrGlobalCf() sql.Connect
		DB() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		//  AllBalanceShardActive return ["partner_code"]orm.BalanceShard
		AllBalanceShardActive() (map[string]orm.BalanceShard, error)
		GetById(id uint32) (orm.BalanceShard, error)
		CreateNew(bll orm.BalanceShard) error
		UpdateAllField(update orm.BalanceShard) error
		UpdateById(id uint32, update map[string]interface{}) error
	}
)

func (rp *BalanceShard) GetById(id uint32) (orm.BalanceShard, error) {
	var balanceShard orm.BalanceShard
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("id = ?", id).First(&balanceShard)
	if rs.Error != nil {
		return balanceShard, rs.Error
	}

	return balanceShard, nil
}

func (rp *BalanceShard) SetConnect(cn sql.Connect) {
	rp.BaseRepo.SetConnect(cn)
}

func (rp BalanceShard) GetConnectFrGlobalCf() sql.Connect {
	//todo get connect from global config
	var cn sql.Connect
	return cn
}
func (rp BalanceShard) DB() sql.Connect {
	return rp.BaseRepo.CN()
}

func (rp *BalanceShard) CreateNew(bll orm.BalanceShard) error {
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

func (rp *BalanceShard) UpdateAllField(update orm.BalanceShard) error {
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

func (rp *BalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Model(&orm.BalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (rp *BalanceShard) AllBalanceShardActive() (map[string]orm.BalanceShard, error) {
	var bsd []orm.BalanceShard
	var rt map[string]orm.BalanceShard

	rp.BaseRepo.UpdateContext()
	if rp.BaseRepo.IsHaveCancelFc() {
		defer rp.BaseRepo.GetCancelFc()
	}

	rs := rp.DB().Where("status = ?", rp.BalanceShardOrm.StatusActive()).Find(&bsd)
	if rs.Error != nil {
		return rt, rs.Error
	}

	for _, v := range bsd {
		rt[v.ShardCode] = v
	}

	return rt, nil
}

//SetTimeout ms
func (rp *BalanceShard) SetTimeout(timeout uint32) {
	rp.BaseRepo.SetTimeout(timeout)
}

//ResetTimeout
func (rp *BalanceShard) ResetTimeout() {
	rp.BaseRepo.ResetTimeout()
}

//SetTimeout ms
func (rp *BalanceShard) SetContext(ctx context.Context) {
	rp.BaseRepo.SetContext(ctx)
}

//ResetTimeout
func (rp *BalanceShard) ResetContext() {
	rp.BaseRepo.ResetContext()
}

func NewBalanceShardRepository() BalanceShardInterface {
	rp := BalanceShard{}
	rp.SetConnect(rp.GetConnectFrGlobalCf())
	return &rp
}

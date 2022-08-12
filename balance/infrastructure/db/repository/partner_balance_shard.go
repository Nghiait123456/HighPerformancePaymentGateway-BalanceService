package repository

import (
	"context"
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type PartnerBalanceShard struct {
	DB         *gorm.DB
	BalanceOrm orm.PartnerBalanceShard
	ctx        context.Context
}

type PartnerBalanceShardInterface interface {
	GetById(id uint32) (orm.PartnerBalanceShard, error)
	CreateNew(bll orm.PartnerBalanceShard) error
	UpdateAllField(update orm.PartnerBalanceShard) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (pbs *PartnerBalanceShard) GetById(id uint32) (orm.PartnerBalanceShard, error) {
	var shard orm.PartnerBalanceShard
	rs := pbs.DB.WithContext(pbs.ctx).Where("id = ?", id).First(&shard)
	if rs.Error != nil {
		return shard, rs.Error
	}

	return shard, nil
}

func (pbs *PartnerBalanceShard) CreateNew(bll orm.PartnerBalanceShard) error {
	rs := pbs.DB.WithContext(pbs.ctx).Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (pbs *PartnerBalanceShard) UpdateAllField(update orm.PartnerBalanceShard) error {
	rs := pbs.DB.WithContext(pbs.ctx).Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (pbs *PartnerBalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rs := pbs.DB.WithContext(pbs.ctx).Model(&orm.PartnerBalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func NewPartnerBalanceShardRepository() PartnerBalanceShardInterface {
	return &PartnerBalanceShard{}
}

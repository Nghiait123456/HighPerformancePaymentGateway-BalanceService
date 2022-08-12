package repository

import (
	"context"
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type BalanceShard struct {
	DB              *gorm.DB
	BalanceShardOrm orm.BalanceShard
	ctx             context.Context
}

type BalanceShardInterface interface {
	GetAllBalanceShard() ([]orm.BalanceShard, error)
	GetById(id uint32) (orm.BalanceShard, error)
	CreateNew(bll orm.BalanceShard) error
	UpdateAllField(update orm.BalanceShard) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (bs *BalanceShard) GetById(id uint32) (orm.BalanceShard, error) {
	var balanceShard orm.BalanceShard
	rs := bs.DB.WithContext(bs.ctx).Where("id = ?", id).First(&balanceShard)
	if rs.Error != nil {
		return balanceShard, rs.Error
	}

	return balanceShard, nil
}

func (bs *BalanceShard) CreateNew(bll orm.BalanceShard) error {
	rs := bs.DB.WithContext(bs.ctx).Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bs *BalanceShard) UpdateAllField(update orm.BalanceShard) error {
	rs := bs.DB.WithContext(bs.ctx).Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bs *BalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rs := bs.DB.WithContext(bs.ctx).Model(&orm.BalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bs *BalanceShard) GetAllBalanceShard() ([]orm.BalanceShard, error) {
	var bsd []orm.BalanceShard
	rs := bs.DB.WithContext(bs.ctx).Find(&bsd)
	if rs.Error != nil {
		return bsd, rs.Error
	}

	return bsd, nil
}

func NewBalanceShardRepository() BalanceShardInterface {
	return &BalanceShard{}
}

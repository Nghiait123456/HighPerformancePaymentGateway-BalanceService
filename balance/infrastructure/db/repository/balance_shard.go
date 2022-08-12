package repository

import (
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type BalanceShard struct {
	DB              *gorm.DB
	BalanceShardOrm orm.BalanceShard
}

type BalanceShardInterface interface {
	GetAllBalanceShard() ([]BalanceShard, error)
	GetById(id uint32) (orm.BalanceShard, error)
	CreateNew(bll orm.BalanceShard) error
	UpdateAllField(update orm.BalanceShard) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (b *BalanceShard) GetById(id uint32) (orm.BalanceShard, error) {
	var balanceShard orm.BalanceShard
	rs := b.DB.Where("id = ?", id).First(&balanceShard)
	if rs.Error != nil {
		return balanceShard, rs.Error
	}

	return balanceShard, nil
}

func (b *BalanceShard) CreateNew(bll orm.BalanceShard) error {
	rs := b.DB.Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceShard) UpdateAllField(update orm.BalanceShard) error {
	rs := b.DB.Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.BalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceShard) GetAllBalanceShard() ([]BalanceShard, error) {
	var bs []BalanceShard
	rs := b.DB.Find(&bs)
	if rs.Error != nil {
		return bs, rs.Error
	}

	return bs, nil
}

func NewBalanceShardRepository() BalanceShardInterface {
	return &BalanceShard{}
}

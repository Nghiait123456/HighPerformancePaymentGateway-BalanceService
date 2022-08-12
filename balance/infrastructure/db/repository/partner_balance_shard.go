package repository

import (
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type PartnerBalanceShard struct {
	DB         *gorm.DB
	BalanceOrm orm.PartnerBalanceShard
}

type PartnerBalanceShardInterface interface {
	GetById(id uint32) (orm.PartnerBalanceShard, error)
	CreateNew(bll orm.PartnerBalanceShard) error
	UpdateAllField(update orm.PartnerBalanceShard) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (b *PartnerBalanceShard) GetById(id uint32) (orm.PartnerBalanceShard, error) {
	var shard orm.PartnerBalanceShard
	rs := b.DB.Where("id = ?", id).First(&shard)
	if rs.Error != nil {
		return shard, rs.Error
	}

	return shard, nil
}

func (b *PartnerBalanceShard) CreateNew(bll orm.PartnerBalanceShard) error {
	rs := b.DB.Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *PartnerBalanceShard) UpdateAllField(update orm.PartnerBalanceShard) error {
	rs := b.DB.Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *PartnerBalanceShard) UpdateById(id uint32, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.PartnerBalanceShard{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func NewPartnerBalanceShardRepository() PartnerBalanceShardInterface {
	return &PartnerBalanceShard{}
}

package repository

import (
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type BalanceLog struct {
	DB         *gorm.DB
	BalanceOrm orm.Balance
}

type BalanceLogInterface interface {
	GetById(id uint32) (orm.BalanceLog, error)
	CreateNew(bll orm.BalanceLog) error
	UpdateAllField(update orm.BalanceLog) error
	UpdateByOrderId(orderId uint64, update map[string]interface{}) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (b *BalanceLog) GetById(id uint32) (orm.BalanceLog, error) {
	var balanceLog orm.BalanceLog
	rs := b.DB.Where("id = ?", id).First(&balanceLog)
	if rs.Error != nil {
		return balanceLog, rs.Error
	}

	return balanceLog, nil
}

func (b *BalanceLog) CreateNew(bll orm.BalanceLog) error {
	rs := b.DB.Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceLog) UpdateAllField(update orm.BalanceLog) error {
	rs := b.DB.Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceLog) UpdateById(id uint32, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.BalanceLog{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *BalanceLog) UpdateByOrderId(orderId uint64, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.BalanceLog{}).Where("order_id = ?", orderId).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func NewBalanceLogRepository() BalanceLogInterface {
	return &BalanceLog{}
}

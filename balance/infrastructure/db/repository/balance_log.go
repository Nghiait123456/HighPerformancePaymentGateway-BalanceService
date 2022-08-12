package repository

import (
	"context"
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type BalanceLog struct {
	DB         *gorm.DB
	BalanceOrm orm.Balance
	ctx        context.Context
}

type BalanceLogInterface interface {
	GetById(id uint32) (orm.BalanceLog, error)
	CreateNew(bll orm.BalanceLog) error
	UpdateAllField(update orm.BalanceLog) error
	UpdateByOrderId(orderId uint64, update map[string]interface{}) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (bl *BalanceLog) GetById(id uint32) (orm.BalanceLog, error) {
	var balanceLog orm.BalanceLog
	rs := bl.DB.WithContext(bl.ctx).Where("id = ?", id).First(&balanceLog)
	if rs.Error != nil {
		return balanceLog, rs.Error
	}

	return balanceLog, nil
}

func (bl *BalanceLog) CreateNew(bll orm.BalanceLog) error {
	rs := bl.DB.WithContext(bl.ctx).Create(&bll)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bl *BalanceLog) UpdateAllField(update orm.BalanceLog) error {
	rs := bl.DB.WithContext(bl.ctx).Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bl *BalanceLog) UpdateById(id uint32, update map[string]interface{}) error {
	rs := bl.DB.WithContext(bl.ctx).Model(&orm.BalanceLog{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (bl *BalanceLog) UpdateByOrderId(orderId uint64, update map[string]interface{}) error {
	rs := bl.DB.WithContext(bl.ctx).Model(&orm.BalanceLog{}).Where("order_id = ?", orderId).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func NewBalanceLogRepository() BalanceLogInterface {
	return &BalanceLog{}
}

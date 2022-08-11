package repository

import (
	"gorm.io/gorm"
	"high-performance-payment-gateway/balance-service/balance/infrastructure/db/orm"
)

type Balance struct {
	DB         *gorm.DB
	BalanceOrm orm.Balance
}

type BalanceInterface interface {
	GetById(id uint32) (orm.Balance, error)
	GetByPartnerCode(partnerCode string) (orm.Balance, error)
	CreateNew(b orm.Balance) error
	UpdateAllField(update orm.Balance) error
	UpdateByPartnerCode(partnerCode string, update map[string]interface{}) error
	UpdateById(id uint32, update map[string]interface{}) error
}

func (b *Balance) GetById(id uint32) (orm.Balance, error) {
	var balance orm.Balance
	rs := b.DB.Where("id = ?", id).First(&balance)
	if rs.Error != nil {
		return balance, rs.Error
	}

	return balance, nil
}

func (b *Balance) GetByPartnerCode(partnerCode string) (orm.Balance, error) {
	var balance orm.Balance
	rs := b.DB.Where("partner_code = ?", partnerCode).First(&balance)
	if rs.Error != nil {
		return balance, rs.Error
	}

	return balance, nil
}

func (b *Balance) CreateNew(bl orm.Balance) error {
	rs := b.DB.Create(&bl)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *Balance) UpdateAllField(update orm.Balance) error {
	rs := b.DB.Updates(&update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *Balance) UpdateByPartnerCode(partnerCode string, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.Balance{}).Where("partner_code = ?", partnerCode).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (b *Balance) UpdateById(id uint32, update map[string]interface{}) error {
	rs := b.DB.Model(&orm.Balance{}).Where("id = ?", id).Updates(update)
	if rs.Error != nil {
		return rs.Error
	}

	return nil
}

func NewBalanceRepository() BalanceInterface {
	return &Balance{}
}
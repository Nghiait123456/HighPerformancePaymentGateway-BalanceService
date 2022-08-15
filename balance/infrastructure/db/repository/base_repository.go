package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type (
	Base struct {
		ctx          context.Context
		CancelFc     context.CancelFunc
		isUseContext bool
		timeOut      uint32 // ms
		isUseTimeout bool
	}

	BaseInterface interface {
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		IsUseTimeout() bool
		IsUseContext() bool
		UpdateContext(db *gorm.DB)
		IsHaveCancelFc() bool
		GetCancelFc() context.CancelFunc
	}
)

//SetTimeout ms
func (b *Base) SetTimeout(timeout uint32) {
	b.isUseTimeout = true
	b.timeOut = timeout
}

//ResetTimeout
func (b *Base) ResetTimeout() {
	b.isUseTimeout = false
	b.timeOut = 0
}

//SetTimeout ms
func (b *Base) SetContext(ctx context.Context) {
	b.isUseContext = true
	b.ctx = ctx
}

//ResetTimeout
func (b *Base) ResetContext() {
	b.isUseContext = false
	b.ctx = nil
}

func (b *Base) IsUseTimeout() bool {
	return b.isUseTimeout == true
}
func (b *Base) IsUseContext() bool {
	return b.isUseContext == true
}

func (b *Base) UpdateContext(db *gorm.DB) {
	var ctx context.Context

	// context is first, use context will disable use timeout
	if b.isUseContext {
		db = db.WithContext(b.ctx)
	} else {
		if b.isUseTimeout {
			ctx, b.CancelFc = context.WithTimeout(context.Background(), time.Duration(b.timeOut)*time.Millisecond)
			db = db.WithContext(ctx)
		}
	}
}

func (b *Base) IsHaveCancelFc() bool {
	return b.CancelFc != nil
}

func (b *Base) GetCancelFc() context.CancelFunc {
	return b.CancelFc
}

func NewBaseRepository() BaseInterface {
	return &Base{}
}

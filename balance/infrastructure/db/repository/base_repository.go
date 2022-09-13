package repository

import (
	"context"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	"time"
)

type (
	Base struct {
		DB           sql.Connect
		ctx          context.Context
		CancelFc     context.CancelFunc
		isUseContext bool
		timeOut      uint32 // ms
		isUseTimeout bool
	}

	BaseInterface interface {
		SetConnect(cn sql.Connect)
		CN() sql.Connect
		SetTimeout(timeout uint32)
		ResetTimeout()
		SetContext(ctx context.Context)
		ResetContext()
		IsUseTimeout() bool
		IsUseContext() bool
		UpdateContext()
		IsHaveCancelFc() bool
		GetCancelFc() context.CancelFunc
	}
)

func (b *Base) SetConnect(cn sql.Connect) {
	b.DB = cn
}

func (b Base) CN() sql.Connect {
	return b.DB
}

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

func (b *Base) UpdateContext() {
	var ctx context.Context

	// context is first, use context will disable use timeout
	if b.isUseContext {
		b.DB = b.DB.WithContext(b.ctx)
	} else {
		if b.isUseTimeout {
			ctx, b.CancelFc = context.WithTimeout(context.Background(), time.Duration(b.timeOut)*time.Millisecond)
			b.DB = b.DB.WithContext(ctx)
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
	rp := Base{}
	return &rp
}

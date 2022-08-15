package shard_balance_logs

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConnectShard struct {
}

type ConnectShardInterface interface {
	ConnectOneShard(dsn string) (Connect, error)
}

func (c ConnectShard) ConnectOneShard(dsn string) (Connect, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func NewConnectShard() ConnectShardInterface {
	return &ConnectShard{}
}

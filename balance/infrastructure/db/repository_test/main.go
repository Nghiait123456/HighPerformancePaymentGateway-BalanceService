package main

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "admin:1adphamnghia@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
		panic(err.Error())
	}

	bl := repository.NewBalanceRepository(db)
	//rs := bl.CreateNew(orm.Balance{
	//	Balance:               1000,
	//	PartnerCode:           "Test",
	//	Status:                "OK",
	//	IndexLogRequestLatest: 1,
	//	CreatedAt:             1,
	//	UpdatedAt:             1,
	//})
	//
	//fmt.Print("rs = ", rs)

	data, err := bl.GetById(1)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
	} else {
		fmt.Print("data", data)
	}

	data2, err2 := bl.GetById(2)
	if err2 != nil {
		fmt.Printf("err: %v", err2.Error())
	} else {
		fmt.Print("data", data2)
	}
}

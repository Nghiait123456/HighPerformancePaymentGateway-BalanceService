package main

import (
	"fmt"
	"github.com/nouney/randomstring"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Balancer struct {
	Id       uint64
	Balancer string
	Status   string
	Detail   string
	CreateAt int
	UpdateAt int
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func FakeOneBalancer() Balancer {
	return Balancer{
		Balancer: "balance_" + randomstring.Generate(4),
		Status:   "status_" + randomstring.Generate(4),
		Detail:   "detail_" + randomstring.Generate(4),
		CreateAt: rand.Int(),
		UpdateAt: rand.Int(),
	}
}

var balancer = []Balancer{}

func UpdateBalancer(db *gorm.DB, id int64, value uint64) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	balancer := Balancer{}

	//Lock the User record with the specified id
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&balancer, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	//update
	err1 := db.Model(&Balancer{}).Where("id = ?", 1).Update("create_at", value)
	if err1.Error != nil {
		panic("insert error" + err1.Error.Error())
	}

	//Commit the transaction and release the lock
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func main() {

	var wg sync.WaitGroup
	username := "admin"
	password := "1adphamnghia"
	dbName := "test_benmark"
	dbHost := "localhost"
	dbPort := "3306"

	dbURI := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4", username, password, dbHost, dbPort, dbName)
	mysqlDrive := mysql.Open(dbURI)
	conn, err := gorm.Open(mysqlDrive)
	if err != nil {
		panic("connect error")
	}

	//rs1 := conn.Raw("SET global max_connections = ?", 100000)
	//if rs1.Error != nil {
	//	panic("set max_connections error" + rs1.Error.Error())
	//}
	var ops uint64

	starTime := time.Now().Unix()
	fmt.Println(starTime)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("start one routine insert, i = ", i)
			for i := 0; i < 50; i++ {
				atomic.AddUint64(&ops, 1)
				err1 := UpdateBalancer(conn, 1, ops)
				if err1 != nil {
					fmt.Println("error", err1)
					panic("insert error" + err1.Error())
				}
			}
		}()
	}

	wg.Wait()
	endTime := time.Now().Unix()

	fmt.Println("rangetime", endTime-starTime, "endtime", endTime, "startTime", starTime)

}

/**
Sumary:

update 50000 row with 1 row:  rangetime 286 endtime 1658585238 startTime 1658584952

The average of mysql is about 250 to 350 qps.
*/

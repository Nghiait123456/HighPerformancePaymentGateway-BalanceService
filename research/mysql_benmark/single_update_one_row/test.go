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
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("start one routine insert, i = ", i)
			for i := 0; i < 100; i++ {
				atomic.AddUint64(&ops, 1)
				err1 := conn.Model(&Balancer{}).Where("id = ?", 1).Update("create_at", ops)
				if err1.Error != nil {
					panic("insert error" + err1.Error.Error())
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

update 50000 row with 1 row:  78 endtime 1658583074 startTime 1658582996
update 150 000 times with 3 row : range time 155 endtime 1658577594 startTime 1658577439

The average of mysql is about 600 to 900 qps.
*/

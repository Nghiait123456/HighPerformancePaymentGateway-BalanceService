package main

import (
	"fmt"
	"github.com/nouney/randomstring"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"sync"
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

	starTime := time.Now().Unix()
	fmt.Println(starTime)
	for i := 0; i < 30000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("start one routine insert, i = ", i)
			for i := 0; i < 100; i++ {
				balancerTest := FakeOneBalancer()
				err1 := conn.Create(balancerTest)
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
i insert 3 M record with time: rangetime 915 endtime 1658575035 startTime 1658574120
On average, mysql executes 3300 insert one second commands.
If using batch insert, this number easily goes up to 300000 command insert in one second.
*/

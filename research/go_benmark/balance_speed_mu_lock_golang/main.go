package main

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type muLockPartner struct {
	lock sync.Mutex
}

var allMutexLockPartner map[int]muLockPartner
var countPartner int

func createAllMutexLock(countPartner int) map[int]muLockPartner {
	if countPartner <= 0 {
		panic("count partner not valid, required greater 0")
	}

	locks := make(map[int]muLockPartner, countPartner)
	for i := 0; i < countPartner; i++ {
		locks[i] = muLockPartner{}
	}

	return locks
}

func init() {
	countPartner = 3000
	allMutexLockPartner = createAllMutexLock(countPartner)

}

func main() {
	app := fiber.New()
	app.Get("/newOrder", func(c *fiber.Ctx) error {
		partnerIdentification := rand.Intn(countPartner)
		mutexOfPartner, ok := allMutexLockPartner[partnerIdentification]
		if !ok {
			err := fmt.Sprintf("mutext for partner  %i doesn't exist", partnerIdentification)
			panic(err)
		}

		mutexOfPartner.lock.Lock()
		//fake time calculator balance in ram
		time.Sleep(4 * time.Microsecond)
		mutexOfPartner.lock.Unlock()

		// fake time insert to DB
		time.Sleep(30 * time.Millisecond)

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	app.Listen(":8080")

}

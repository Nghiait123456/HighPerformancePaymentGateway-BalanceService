package main

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"time"
)

var maxWorker = 30000

type balancerWorker struct {
	amount      int
	partnerCode int
	timeout     <-chan time.Time
	resProcess  chan resFrProcessBalance
}

type resFrProcessBalance struct {
	status       bool
	errorCode    uint
	errorMessage string
	httpStatus   int
}

func main() {
	timeout := time.After(5 * time.Second)
	resFrProcess := make(chan resFrProcessBalance)
	channelWorker := make(chan balancerWorker, 2)
	createOneWorkerAndHandleJob(channelWorker, 1)
	fmt.Println("init worker success")
	app := fiber.New()

	app.Get("/testAddJob", func(c *fiber.Ctx) error {
		requestInfo := balancerWorker{amount: 1, partnerCode: 1, timeout: timeout, resProcess: resFrProcess}
		channelWorker <- requestInfo
		return c.Status(200).JSON(fiber.Map{
			"status": "add job success",
		})
	})

	app.Listen(":8082")

}

func createOneWorkerAndHandleJob(b chan balancerWorker, partnerCode int) {
	go func() {
		fmt.Println("handle worker:  create loop forever")
		for {
			fmt.Println("11111")
			select {
			case workInfo := <-b:
				fmt.Println("start job, amount %i", workInfo.amount)
				time.Sleep(20 * time.Second)
				fmt.Println("end job, amount %i", workInfo.amount)
				//new Chan
				//go func() {
				//
				//}()
				//fmt.Println("start run workInfo")
				//if workInfo.partnerCode != partnerCode {
				//	err := fmt.Sprintf("put job to wrong channel, partnerCode right %i,  partnerCode pass:  %i", partnerCode, workInfo.partnerCode)
				//	panic(err)
				//}
				//
				//fmt.Println("run default case ")
				////atomic.AddUint64(&up, 1)
				//break
				//}
			}
		}
	}()

}

package main

import (
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

type balancer struct {
	amount      int
	partnerCode int
	resProcess  chan resFrProcessBalance
}

type balancerWorker struct {
	amount      int
	partnerCode int
	resProcess  chan resFrProcessBalance
}

type resFrProcessBalance struct {
	status       bool
	errorCode    uint
	errorMessage string
	httpStatus   int
}

var up uint64
var down uint64
var maxGroupPartner = 30000

type JsonResponse map[string]interface{}

func main() {
	startTime := time.Now().UnixNano()
	down = 1000000000
	up = 0
	inputRequest := make(chan balancer, 1000000)

	fmt.Println("start create multi worker")
	works, err := createMultiWorker()
	if err != nil {
		panic("error createMultiWorker, error code: " + err.Error())
	}
	fmt.Println("end create multi worker")

	//fmt.Println("start push job")
	//Push(inputRequest)
	//fmt.Println("end push job")

	fmt.Println("start LB ")
	loadBalancerJobWorker(inputRequest, works)

	fmt.Println("init fiber")
	app := fiber.New()

	app.Get("/getDown", func(c *fiber.Ctx) error {
		c.Status(http.StatusOK)
		return c.JSON(fiber.Map{
			"up":                     "up",
			"time":                   time.Now().UnixNano() - startTime,
			"average message/second": up*10 ^ 9/(uint64(time.Now().UnixNano()-startTime)),
		})
	})

	app.Get("/newOrder", func(c *fiber.Ctx) error {
		timeoutProcess := time.After(800 * time.Millisecond)
		resFrProcess := make(chan resFrProcessBalance)

		b := balancer{amount: 1, partnerCode: rand.Intn(maxGroupPartner), resProcess: resFrProcess}
		fmt.Println("start add job to channel")
		inputRequest <- b
		fmt.Println("done add job to channel")

		for {
			select {
			case resP := <-resFrProcess:
				return c.Status(resP.httpStatus).JSON(fiber.Map{
					"message": resP.errorMessage,
					"detail":  "response form process",
				})

			case <-timeoutProcess:
				fmt.Println("Timeout request")
				c.Status(http.StatusSeeOther)
				return c.JSON(JsonResponse{
					"status": "timeout",
					"detail": "response form timeout request",
				})
			}
		}
	})

	app.Get("/testTimeSleep", func(c *fiber.Ctx) error {
		fmt.Println("startSleep")
		time.Sleep(10 * time.Second)
		fmt.Println("endSleep")
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	app.Listen(":8081")
}

func loadBalancerJobWorker(inputJob chan balancer, channelsWorker map[int]chan balancerWorker) {
	for i := 0; i < 150000; i++ {
		go func() {
			for {
				select {
				// move job to worker handle it
				case workInfo := <-inputJob:
					if channelOneWorker, ok := channelsWorker[workInfo.partnerCode]; ok {
						oneJobInfo := balancerWorker{amount: workInfo.amount, partnerCode: workInfo.partnerCode, resProcess: workInfo.resProcess}
						channelOneWorker <- oneJobInfo
					} else {
						err := fmt.Sprintf("job over range in worker %i", workInfo.partnerCode)
						panic(err)
					}
				}
			}
		}()
	}

}

func createMultiWorker() (map[int]chan balancerWorker, error) {
	b := make(map[int]chan balancerWorker)
	for i := 0; i < maxGroupPartner; i++ {
		b[i] = make(chan balancerWorker)
		createOneWorkerAndHandleJob(b[i], i)
	}

	return b, nil
}

func createOneWorkerAndHandleJob(b chan balancerWorker, partnerCode int) {
	go func() {
		for {
			select {
			case workInfo := <-b:
				if workInfo.partnerCode != partnerCode {
					err := fmt.Sprintf("put job to wrong channel, partnerCode right %i,  partnerCode pass:  %i", partnerCode, workInfo.partnerCode)
					panic(err)
				}
				atomic.AddUint64(&up, 1)
				res := resFrProcessBalance{status: true, errorCode: 2, errorMessage: "success job", httpStatus: http.StatusOK}
				//fake timeout process calculator and save DB
				time.Sleep(300 * time.Millisecond)
				workInfo.resProcess <- res
				fmt.Println("done one job")
			}
		}
	}()

}

func Push(c chan balancer) {
	for i := 0; i < 150000; i++ {
		go func() {
			for i := 0; i < maxGroupPartner; i++ {
				b := balancer{amount: 1, partnerCode: i}
				c <- b
			}
		}()
	}
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
	"time"
)

type balancer struct {
	amount      int
	partnerCode int
}

type balancerWorker struct {
	amount      int
	partnerCode int
}

var up uint64
var down uint64
var maxWorker = 30000

func main() {
	startTime := time.Now().UnixNano()
	down = 1000000000
	up = 0
	inputRequest := make(chan balancer)

	fmt.Println("start create multi worker")
	works, err := createMultiWorker()
	if err != nil {
		panic("error createMultiWorker, error code: " + err.Error())
	}
	fmt.Println("end create multi worker")

	fmt.Println("start push job")
	Push(inputRequest)
	fmt.Println("end push job")

	fmt.Println("start LB ")
	loadBalancerJobWorker(inputRequest, works)

	go func() {
		var limit = uint64(maxWorker) * 50000
		for {
			if up == limit {
				fmt.Printf("end time run job, number request %i, rangeTime %i", limit, time.Now().UnixNano()-startTime)
				return
			}
		}
	}()

	fmt.Println("init gin")
	r := gin.Default()
	r.GET("/getDown", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"down": down,
		})
	})

	r.GET("/getUp", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"up":                     up,
			"time":                   time.Now().UnixNano() - startTime,
			"average message/second": up*10 ^ 9/(uint64(time.Now().UnixNano()-startTime)),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func loadBalancerJobWorker(inputJob chan balancer, channelsWorker map[int]chan balancerWorker) {
	for i := 0; i < 150000; i++ {
		go func() {
			for {
				select {
				// move job to worker handle it
				case workInfo := <-inputJob:
					if channelOneWorker, ok := channelsWorker[workInfo.partnerCode]; ok {
						oneJobInfo := balancerWorker{amount: workInfo.amount, partnerCode: workInfo.partnerCode}
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
	for i := 0; i < maxWorker; i++ {
		b[i] = make(chan balancerWorker)
		createOneWorkerAndChannel(b[i], i)
	}

	return b, nil
}

func createOneWorkerAndChannel(b chan balancerWorker, partnerCode int) {
	go func() {
		for {
			select {
			case workInfo := <-b:
				if workInfo.partnerCode != partnerCode {
					err := fmt.Sprintf("put job to wrong channel, partnerCode right %i,  partnerCode pass:  %i", partnerCode, workInfo.partnerCode)
					panic(err)
				}

				atomic.AddUint64(&up, 1)
			}
		}
	}()

}

func Push(c chan balancer) {
	for i := 0; i < 150000; i++ {
		go func() {
			for i := 0; i < maxWorker; i++ {
				b := balancer{amount: 1, partnerCode: i}
				c <- b
			}
		}()
	}
}

/**
end time run job, number request %!i(uint64=1500000000), rangeTime %!i(int64=1193853098263)

I have demo job send request check balance, handle in local in memory with golang. I use channel golang.In my laptop  4 core, 2.2 Ghz, i75600, i get 1.5M to 2 M request transfer in second.

If only talk about channel transfer message, channel can convert 14M to 20M message in second. This depends in part on the machine configuration.
*/

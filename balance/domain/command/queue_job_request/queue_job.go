package queue_job_request

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	OneRequest = calculator.BalancerRequest
	JobRequest chan OneRequest

	QueueJob struct {
		QJob       JobRequest
		AllPartner calculator.AllPartnerInterface
	}

	QueueJobInterface interface {
		Push(rq OneRequest)
		AutoHandleRequest()
		Init(allP calculator.AllPartnerInterface)
	}
)

func (q *QueueJob) Push(rq OneRequest) {
	q.QJob <- rq
	fmt.Println("push job done")
}

func (q *QueueJob) AutoHandleRequest() {
	log.Info("start auto handle request")
	for {
		select {
		case rq := <-q.QJob:
			pn, errGOP := q.AllPartner.GetOnePartner(rq.PartnerCode)
			if errGOP != nil {
				log.WithFields(log.Fields{
					"errMessage": errGOP.Error(),
				}).Error("dont get one partner for handle request")

				//todo push error, push even, log
				break
			}

			rs, errHOR := pn.HandleOneRequestBalance(rq)
			if rs != true {
				log.WithFields(log.Fields{
					"errMessage": errHOR.Error(),
				}).Error("dont get one partner for handle request")
				//todo push error, push even, log
				break
			}

			q.AllPartner.UpdateOnePartner(pn)
			// todo push error, push even, log success
		}
	}
}

func (q *QueueJob) Init(allP calculator.AllPartnerInterface) {
	q.AllPartner = allP
	go func() {
		q.AutoHandleRequest()
	}()

	rs := q.testInitQueue()
	if rs != true {
		m := "Queue Job handle request balance init error"
		log.Error(m)
		fmt.Println(m)
		q.AllPartner.ThrowEStop()
		//todo alert message warring
	}
}

func (q *QueueJob) testInitQueue() bool {
	//sleep for make sure consumer ready
	time.Sleep(100 * time.Nanosecond)
	tOut := time.After(1 * time.Millisecond)
	oneRqTest := OneRequest{
		AmountRequest:         0,
		PartnerCode:           "dev_test_channel_ready_123456",
		PartnerIdentification: 0,
		OrderID:               0,
		TypeRequest:           "test",
	}

	select {
	case q.QJob <- oneRqTest:
		message := fmt.Sprintf("testInitQueue: Success, queue ready for work")
		log.Info(message)
		fmt.Println(message)
		return true
	case <-tOut:
		message := fmt.Sprintf("testInitQueue: Error, queue dont have consumer or queue work very late")
		log.Error(message)
		fmt.Println(message)
		return false
	}
}

func NewQueueJob() QueueJobInterface {
	return &QueueJob{
		QJob: make(JobRequest),
	}
}

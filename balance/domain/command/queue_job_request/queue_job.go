package queue_job_request

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"
	log "github.com/sirupsen/logrus"
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
}

func (q *QueueJob) AutoHandleRequest() {
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
				// todo push error, push even, log
				break
			}

			// todo push error, push even, log success
		}
	}
}

func (q *QueueJob) Init(allP calculator.AllPartnerInterface) {
	q.AllPartner = allP
	go func() {
		q.AutoHandleRequest()
	}()
}

func NewQueueJob() QueueJobInterface {
	return &QueueJob{}
}

package queue_job_request

import "github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"

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
	}
)

func (q *QueueJob) Push(rq OneRequest) {
	q.QJob <- rq
}

func (q *QueueJob) AutoHandleRequest() {
	for {
		select {
		case rq := <-q.QJob:
			pn, errPN := q.AllPartner.GetOnePartner(rq.PartnerCode)
			if errPN != nil {
				// todo push error, push even, log
			}

			rs, _ := pn.HandleOneRequestBalance(rq)
			if rs != true {
				// todo push error, push even, log
			}
		}
	}

}

func NewQueueJob() QueueJobInterface {
	return &QueueJob{}
}

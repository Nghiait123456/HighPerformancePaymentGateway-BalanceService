package application

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/application/respone_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/calculator"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/queue_job_request"
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	AllPartnerBalance struct {
		AllPartner calculator.AllPartnerInterface
		QueueJob   queue_job_request.QueueJobInterface
	}

	BalanceRequest struct {
		BRequest calculator.BalancerRequest
	}

	AllPartnerBalanceInterface interface {
		Init() error
		HandleRequestBalance(b BalanceRequest) (bool, respone_request_balance.RequestBalanceResponse)
	}
)

func (a *AllPartnerBalance) Init() error {
	err := a.AllPartner.InitAllPartnerInfo()
	if err != nil {
		log.WithFields(log.Fields{
			"errMessage": err.Error(),
		}).Error("InitAllPartnerInfo")

		panic(fmt.Sprintf("InitAllPartnerInfo error: %s", err.Error()))
		os.Exit(0)
	}

	a.QueueJob.Init(a.AllPartner)
	return nil
}

func (a *AllPartnerBalance) HandleRequestBalance(b BalanceRequest) (bool, respone_request_balance.RequestBalanceResponse) {
	a.QueueJob.Push(b.BRequest)
	return true, respone_request_balance.SuccessBalanceResponse()
}

func NewAllPartnerBalance(allPartner calculator.AllPartnerInterface, queueJob queue_job_request.QueueJobInterface) AllPartnerBalanceInterface {
	return &AllPartnerBalance{
		AllPartner: allPartner,
		QueueJob:   queueJob,
	}
}

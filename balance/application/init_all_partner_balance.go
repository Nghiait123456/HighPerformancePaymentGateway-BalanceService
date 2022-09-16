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

	GroupBalanceRequest struct {
		BRequest           BGroupRequest
		ListRequestSuccess ListRequestSuccess
	}

	BGroupRequest            = []calculator.BalancerRequest
	ListRequestSuccess       = []calculator.BalancerRequest
	DetailResultGroupRequest = GroupBalanceRequest

	AllPartnerBalanceInterface interface {
		Init() error
		HandleOneRequestBalance(b BalanceRequest) (respone_request_balance.RequestBalanceResponse, bool)
		HandleGroupRequestBalance(gb GroupBalanceRequest) (respone_request_balance.RequestBalanceResponse, DetailResultGroupRequest, bool)
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

func (a *AllPartnerBalance) HandleOneRequestBalance(b BalanceRequest) (respone_request_balance.RequestBalanceResponse, bool) {
	if a.AllPartner.IsEStop() {
		return respone_request_balance.ErrorWhySystemEStop(), false
	}

	a.QueueJob.Push(b.BRequest)
	return respone_request_balance.SuccessBalanceResponse(), true
}

func (a *AllPartnerBalance) HandleGroupRequestBalance(gb GroupBalanceRequest) (respone_request_balance.RequestBalanceResponse, DetailResultGroupRequest, bool) {
	if a.AllPartner.IsEStop() {
		return respone_request_balance.ErrorWhySystemEStop(), DetailResultGroupRequest{}, false
	}

	for _, v := range gb.BRequest {
		a.QueueJob.Push(v)
		gb.ListRequestSuccess = append(gb.ListRequestSuccess, v)
	}

	detail := gb
	return respone_request_balance.SuccessBalanceResponse(), detail, true
}

func NewAllPartnerBalance(allPartner calculator.AllPartnerInterface, queueJob queue_job_request.QueueJobInterface) *AllPartnerBalance {
	var _ AllPartnerBalanceInterface = (*AllPartnerBalance)(nil)
	return &AllPartnerBalance{
		AllPartner: allPartner,
		QueueJob:   queueJob,
	}
}

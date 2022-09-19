package calculator

import (
	"errors"
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/domain/command/logs_request_balance"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/db/connect/sql"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

/**
all information partner for calculator balance
*/

type (
	AllPartner struct {
		allPartner        PartnersBalance
		muLock            sync.Mutex
		cnRechargeLog     CnRechargeLog
		cnBalance         CnBalance
		logRequestBalance logs_request_balance.LogsInterface
		eStop             emergencyStopInterface
	}

	CnRechargeLog   sql.Connect
	CnBalance       sql.Connect
	PartnersBalance map[string]partnerBalance

	AllPartnerInterface interface {
		initPartnersInterface
	}

	initPartnersInterface interface {
		LoadAllPartnerInfo() (PartnersBalance, error)
		InitAllPartnerInfo() error
		UpdateOnePartner(p partnerBalance) error
		GetOnePartner(partnerCode string) (partnerBalance, error)
		getKeyOnePartner(p partnerBalance) string
		ThrowEStop()
		IsEStop() bool
	}
)

func (allP *AllPartner) LoadAllPartnerInfo() (PartnersBalance, error) {
	fake := make(PartnersBalance)
	//todo get indexLogRequestLatest from DB and update to
	fake["TEST"] = partnerBalance{
		partnerCode:           "TEST",
		partnerName:           "TEST",
		partnerIdentification: 1,
		balance:               99999000000000,
		amountPlaceHolder:     0,
		cnRechargeLog:         allP.cnRechargeLog,
		cnBalance:             allP.cnBalance,
		logRequestBalance:     allP.logRequestBalance,
	}

	fake["TEST1"] = partnerBalance{
		partnerCode:           "TEST1",
		partnerName:           "TEST1",
		partnerIdentification: 1,
		balance:               99999000000000,
		amountPlaceHolder:     0,
		cnRechargeLog:         allP.cnRechargeLog,
		cnBalance:             allP.cnBalance,
		logRequestBalance:     allP.logRequestBalance,
	}

	return fake, nil
}

func (allP *AllPartner) InitAllPartnerInfo() error {
	allPartner, err := allP.LoadAllPartnerInfo()
	if err != nil {
		return err
	}

	balancePlaceHolderHistory := NewBalancePlaceHolderHistory()

	allP.muLock.Lock()
	//init raw AllPartner
	for k, v := range allPartner {
		allP.allPartner[k] = v
	}

	//merger balancePlaceHolderHistory to AllPartner
	for k, v := range allP.allPartner {
		placeHolder, ok := balancePlaceHolderHistory.GetAllPlaceHolder()[k]
		if ok {
			v.amountPlaceHolder = placeHolder.amountPlaceHolder
			allP.allPartner[k] = v
		} else {
			v.amountPlaceHolder = 0
		}
	}

	allP.muLock.Unlock()
	//allP.dumpAllPartnerInfo()

	return nil
}

func (allP *AllPartner) dumpAllPartnerInfo() {
	go func() {
		for {
			startTimeShow := "-------------------------------- start times show -----------------------------------------------------------------------------"
			fmt.Println(startTimeShow)
			log.Info(startTimeShow)

			for _, v := range allP.allPartner {
				show := fmt.Sprintf("partnerCode: %s, amount: %d , amountPlaceHolder: %d", v.partnerCode, v.balance, v.amountPlaceHolder)
				fmt.Println(show)
				log.Info(show)
			}

			endTimeShow := "############################################### end times show ###################################################################"
			fmt.Println(endTimeShow)
			log.Info(endTimeShow)
			time.Sleep(5000 * time.Millisecond)
		}
	}()
}

func (allP *AllPartner) UpdateOnePartner(p partnerBalance) error {
	key := allP.getKeyOnePartner(p)
	allP.allPartner[key] = p

	return nil
}

func (allP *AllPartner) getKeyOnePartner(p partnerBalance) string {
	return p.partnerCode
}

func (allP *AllPartner) GetOnePartner(partnerCode string) (partnerBalance, error) {
	partner, ok := allP.allPartner[partnerCode]
	if !ok {
		err := fmt.Sprintf("partnercode %s not exists", partnerCode)
		return partnerBalance{}, errors.New(err)
	}

	return partner, nil
}

func (allP *AllPartner) ThrowEStop() {
	allP.eStop.ThrowEmergencyStop()
}

func (allP AllPartner) IsEStop() bool {
	return allP.eStop.IsStop()
}

func InitAllPartnerData() AllPartnerInterface {
	allPartner := AllPartner{}
	err := allPartner.InitAllPartnerInfo()
	if err != nil {
		panic("Init all partner error: " + err.Error())
	}

	return &allPartner
}

func NewAllPartner(allPartner PartnersBalance, cnRechargeLog CnRechargeLog, cnBalance CnBalance, logRequestBalance logs_request_balance.LogsInterface) *AllPartner {
	var _ AllPartnerInterface = (*AllPartner)(nil)
	return &AllPartner{
		allPartner:        allPartner,
		cnRechargeLog:     cnRechargeLog,
		cnBalance:         cnBalance,
		logRequestBalance: logRequestBalance,
		eStop:             NewEmergencyStop(),
	}
}

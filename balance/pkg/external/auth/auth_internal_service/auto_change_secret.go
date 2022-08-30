package auth_internal_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type (
	AutoChangeSecret struct {
		timeUpdateSecret uint64 // ms
	}

	AutoChangeSecretInterface interface {
		Init() error
		getSecretFrRemote() (AllSecret, error)
		replaceNewSecretInRemote() error
		createNewSecret(currentSecret string) (string, error)
		autoUpdateNewTimeChangeSecret() error
		checkAllConfigEnvBeforeRun() error
	}
)

const (
	IS_AUTO_CHANGE_SECRET_REMOTE_KEY       = "IS_AUTO_CHANGE_SECRET_REMOTE_KEY"
	IS_AUTO_CHANGE_SECRET_REMOTE_VALUE     = "true"
	IS_NOT_AUTO_CHANGE_SECRET_REMOTE_VALUE = "false"
)

func (a *AutoChangeSecret) getSecretFrRemote() (AllSecret, error) {
	allS := AllSecret{}
	sAws := env.NewAwsManagerSecret()
	secretName := os.Getenv(SECRET_NAME_KEY)
	if secretName == "" {
		log.WithFields(log.Fields{
			"secretName": "",
		}).Error("Secret name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	region := os.Getenv(REGION_KEY)
	if region == "" {
		log.WithFields(log.Fields{
			"region": region,
		}).Error("region name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	versionState, errVST := os.LookupEnv(VERSION_STATE_KEY)
	if errVST != true {
		log.WithFields(log.Fields{
			"versionState": "",
		}).Error("versionState dont have config in env")
		panic("versionState dont have config in env")
		os.Exit(0)
	}
	sAws.Init(secretName, region, versionState)

	sString, errGetSec := sAws.GetSecret()
	if errGetSec != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetSec.Error(),
		}).Error("get secret from remote auth_internal's service error")
		return allS, errGetSec
	}

	errDecode := json.Unmarshal([]byte(sString), &allS)
	if errDecode != nil {
		log.WithFields(log.Fields{
			"errorMessage": errDecode.Error(),
		}).Error("secret save in remote wrong format")
		return allS, errGetSec
	}

	if allS.SecretCurrent == "" {
		return allS, NewErrorDontConstructSecretFirstTime()
	}

	return allS, nil
}

func (a AutoChangeSecret) checkAllConfigEnvBeforeRun() error {
	if os.Getenv(SECRET_NAME_KEY) == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", SECRET_NAME_KEY))
	}

	if os.Getenv(REGION_KEY) == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", REGION_KEY))
	}

	_, e := os.LookupEnv(VERSION_STATE_KEY)
	if e != true {
		return errors.New(fmt.Sprintf("missing config env key %s", VERSION_STATE_KEY))
	}

	if os.Getenv(IS_AUTO_CHANGE_SECRET_REMOTE_KEY) == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", IS_AUTO_CHANGE_SECRET_REMOTE_KEY))
	}

	return nil

}

func (a *AutoChangeSecret) replaceNewSecretInRemote() error {
	s, errGetFrRemote := a.getSecretFrRemote()
	if errGetFrRemote != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetFrRemote.Error(),
		}).Error("getSecretFrRemote have error")
		return errGetFrRemote
	}

	newSecret, errNS := a.createNewSecret(s.SecretCurrent)
	if errNS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errNS.Error(),
		}).Panicf("createNewSecret have error")
		return errNS
	}

	secretU := AllSecret{
		SecretCurrent:    newSecret,
		SecretNearestOld: s.SecretCurrent,
		TimeUpdateSecret: s.TimeUpdateSecret,
	}

	secretUB, errCS := json.Marshal(secretU)
	if errCS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errCS.Error(),
		}).Panicf("convert secret update error")
		return errCS
	}
	secretUString := string(secretUB)

	secretName := os.Getenv(SECRET_NAME_KEY)
	if secretName == "" {
		log.WithFields(log.Fields{
			"secretName": "",
		}).Error("Secret name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	region := os.Getenv(REGION_KEY)
	if region == "" {
		log.WithFields(log.Fields{
			region: "",
		}).Error("region name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	versionState, errVST := os.LookupEnv(VERSION_STATE_KEY)
	if errVST != true {
		log.WithFields(log.Fields{
			"versionState": "",
		}).Error("versionState dont have config in env")
		panic("versionState dont have config in env")
		os.Exit(0)
	}

	sAws := env.NewAwsManagerSecret()
	sAws.Init(secretName, region, versionState)
	errUSS := sAws.UpdateSecretString(secretUString)
	if errUSS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errUSS.Error(),
		}).Error("update secret error")
	}

	log.Infof("ReplaceNewSecretInRemote success %s", secretUString)
	return nil
}

func (a AutoChangeSecret) createNewSecret(currentSecret string) (string, error) {
	version, errV := a.getVersionFrSecret(currentSecret)
	if errV != nil {
		log.WithFields(log.Fields{
			"errorMessage": errV.Error(),
		}).Error("format secret have error")
		panic(fmt.Sprintf("format secret_current from auth_internal service have error: %s", errV.Error()))
	}

	version++
	versionString := a.standardizedVersion(version)
	bodyString := uniuri.NewLen(LENGTH_BODY_STRING)

	secretNew := bodyString + versionString
	return secretNew, nil
}

func (a AutoChangeSecret) standardizedVersion(v uint64) string {
	vNew := strconv.FormatUint(v, 10)
	nVewLen := len(vNew)

	if nVewLen < LENGTH_VERSION_STRING {
		for i := 0; i < LENGTH_VERSION_STRING-nVewLen; i++ {
			vNew = "0" + vNew
		}
	}

	return vNew
}

func (a *AutoChangeSecret) autoUpdateNewTimeChangeSecret() error {
	allS, errGS := a.getSecretFrRemote()
	if errGS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGS.Error(),
		}).Error("dont get secret from remote")
	}

	timeU, errParserTime := strconv.ParseUint(allS.TimeUpdateSecret, 10, 64)
	if errParserTime != nil {
		log.WithFields(log.Fields{
			"errorMessage": errParserTime.Error(),
		}).Error("format secret have error")
		return errParserTime
	}

	a.timeUpdateSecret = timeU
	return nil
}

func (a AutoChangeSecret) Init() error {
	autoChangeRemoteSecretMode := os.Getenv(IS_AUTO_CHANGE_SECRET_REMOTE_KEY)
	if autoChangeRemoteSecretMode == IS_AUTO_CHANGE_SECRET_REMOTE_VALUE {
		go func() {
			a.autoReplaceNewSecretInRemote()
		}()
	}

	return nil
}

func (a AutoChangeSecret) autoReplaceNewSecretInRemote() {
	for {
		errT := a.autoUpdateNewTimeChangeSecret()
		if errT != nil {
			log.WithFields(log.Fields{
				"errorMessage": errT.Error(),
			}).Error("autoUpdateNewTimeChangeSecret error")
		} else {
			log.Info("autoUpdateNewTimeChangeSecret success")
		}

		errRNS := a.replaceNewSecretInRemote()
		if errRNS != nil {
			log.WithFields(log.Fields{
				"errorMessage": errRNS.Error(),
			}).Error("replace secret error")
		} else {
			log.Info("autoReplaceNewSecretInRemote success")
		}

		// Sleep() is the best for optimal cpu resource
		time.Sleep(time.Duration(a.timeUpdateSecret) * time.Millisecond)
	}
}

func (a AutoChangeSecret) validateSecret(secret string) (bool, error) {
	if len(secret) != LENGTH_SECRET_STRING {
		return false, errors.New(fmt.Sprintf("len secret must is %d", LENGTH_SECRET_STRING))
	}

	return true, nil
}

func (a AutoChangeSecret) getVersionFrSecret(secret string) (uint64, error) {
	if len(secret) != LENGTH_SECRET_STRING {
		return 0, errors.New(fmt.Sprintf("length secret must is %d", LENGTH_SECRET_STRING))
	}

	v, err := strconv.ParseUint(secret[20:40], 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("version secret %v is not uint, wrong format", secret[20:40]))
	}

	return v, nil
}

// run when first time init
func FirstTimeContructSecret() error {
	allS := AllSecret{}
	allS.SecretCurrent = createFirstSecret()
	allS.SecretNearestOld = createFirstSecret()
	allS.TimeUpdateSecret = strconv.Itoa(TIME_UPDATE_SECRET_DEFAULT)

	secretUB, errCS := json.Marshal(allS)
	if errCS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errCS.Error(),
		}).Error("convert secret update error")
		return errCS
	}
	secretUString := string(secretUB)

	secretName := os.Getenv(SECRET_NAME_KEY)
	if secretName == "" {
		log.WithFields(log.Fields{
			"secretName": "",
		}).Error("Secret name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	region := os.Getenv(REGION_KEY)
	if region == "" {
		log.WithFields(log.Fields{
			region: "",
		}).Error("region name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	versionState, errVST := os.LookupEnv(VERSION_STATE_KEY)
	if errVST != true {
		log.WithFields(log.Fields{
			"versionState": "",
		}).Error("versionState dont have config in env")
		panic("versionState dont have config in env")
		os.Exit(0)
	}

	sAws := env.NewAwsManagerSecret()
	sAws.Init(secretName, region, versionState)
	errUSS := sAws.UpdateSecretString(secretUString)
	if errUSS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errUSS.Error(),
		}).Error("update secret error")
		return errUSS
	}

	return nil
}

func (a AutoChangeSecret) ListEnvRequireSetupBeforeRunPacket() []string {
	return []string{
		SECRET_NAME_KEY,
		REGION_KEY,
		VERSION_STATE_KEY,
		IS_AUTO_CHANGE_SECRET_REMOTE_KEY,
	}
}

func NewAutoChangeSecret() AutoChangeSecretInterface {
	return &AutoChangeSecret{}
}

package auth_internal_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"time"
)

/**
  auth internal in microservice: high performance + security + handle all race conditions + zero down time when update new secret
  process:
                                      you have:    secret manager remote(encry AllData) save all secret key
   all service ==> cal to secret manager remote ==> get secret ==> auth   |||||||||||||||||||||||||||  secret manager remote: auto change secret for time

*/
type (
	AuthInternalService struct {
		secretCurrent     string
		secretNearestOld  string
		versionCurrent    uint64
		versionNearestOld uint64
		timeUpdateSecret  uint64 // ms
	}

	AuthResponse struct {
		HttpCode int
		Status   string
		Code     int
		Message  string
	}

	AllSecret struct {
		SecretNearestOld string `json:"SECRET_NEAREST_OLD"`
		SecretCurrent    string `json:"SECRET_CURRENT"`
		TimeUpdateSecret string `json:"TIME_UPDATE_SECRET"`
	}

	AuthInternalServiceInterface interface {
		/*
		  list fc for auth
		*/
		Auth(secretCheck string) (bool, AuthResponse) // auth to incoming request
		InitAuth() error
		autoUpdateSecretInLocal()
		UpdateSecretInLocal() error
		getSecretFrRemote() (AllSecret, error)
		checkAllConfigEnvBeforeRun() error
		ListEnvRequireSetupBeforeRunPacket() []string
	}
)

const (
	LENGTH_SECRET_STRING       = 40
	LENGTH_VERSION_STRING      = 20
	LENGTH_BODY_STRING         = 20
	TIME_UPDATE_SECRET_DEFAULT = 10000 //ms

	//all value map key of secret require pass when init data
	SECRET_NAME_KEY                = "AUTH_INTERNAL_SECRET_NAME"
	REGION_KEY                     = "AUTH_INTERNAL_VERSION"
	VERSION_STATE_KEY              = "AUTH_INTERNAL_VERSION_STATE"
	IS_USE_AUTH_INTERNAL_KEY       = "IS_USE_AUTH_INTERNAL_KEY"
	IS_USE_AUTH_INTERNAL_VALUE     = "true"
	IS_NOT_USE_AUTH_INTERNAL_VALUE = "false"
)

// Auth one of the parties has not updated the latest secret, will use the nearest secret to auth. ==> zero down time why race conditions update secret
func (a *AuthInternalService) Auth(secretCheck string) (bool, AuthResponse) {
	if len(secretCheck) != LENGTH_SECRET_STRING {
		log.WithFields(log.Fields{
			"secretCheck": secretCheck,
		}).Errorf("length SecretChek must is %d", LENGTH_SECRET_STRING)
		return responseAuthErrorWhyVersionSecretKeyWrongFormat()
	}

	/**
	  check with current secret
	*/
	if a.secretCurrent == secretCheck {
		return responseAuthSuccess()
	}

	vSecretCheck, errV := a.getVersionFrSecret(secretCheck)
	if errV != nil {
		log.WithFields(log.Fields{
			"version":      vSecretCheck,
			"errorMessage": errV.Error(),
		}).Error("Wrong format version secret key input check auth")
		return responseAuthErrorWhyVersionSecretKeyWrongFormat()
	}

	// if have new version, update secret and check again
	if vSecretCheck > a.versionCurrent {
		err := a.UpdateSecretInLocal()
		if err != nil {
			log.WithFields(log.Fields{
				"errMessage": err.Error(),
			}).Error("Update secret error")
			return defaultErrorWhyNotClearReason()
		}

		if vSecretCheck != a.versionCurrent {
			return responseAuthErrorWhyWrongSecret()
		}

		if secretCheck == a.secretCurrent {
			return responseAuthSuccess()
		} else {
			return responseAuthErrorWhyWrongSecret()
		}
	}

	if vSecretCheck == a.versionCurrent {
		return responseAuthErrorWhyWrongSecret()
	}

	/*
	  check with secretNearestOld, in case request service not in time update new secret, they still auth success with secretNearestOld
	*/
	if secretCheck == a.secretNearestOld {
		return responseAuthSuccessButUpdateNewSecret()
	}

	if vSecretCheck < a.versionNearestOld {
		log.WithFields(log.Fields{
			"versionInput":      vSecretCheck,
			"versionNearestOld": a.versionNearestOld,
		}).Error("Version of secret is old, please update and try again")

		return responseAuthErrorWhyVersionSecretOld()
	}

	if vSecretCheck == a.versionNearestOld {
		log.WithFields(log.Fields{
			"versionInput":      vSecretCheck,
			"versionNearestOld": a.versionNearestOld,
		}).Error("Version of secret is old, please update and try again")

		return responseAuthErrorWhyWrongSecret()
	}

	return defaultErrorWhyNotClearReason()
}

func (a AuthInternalService) validateAllSecret() (bool, error) {
	if len(a.secretCurrent) != 40 {
		return false, errors.New("secret wrong format")
	}

	if len(a.secretNearestOld) != 40 {
		return false, errors.New("secret wrong format")
	}

	return true, nil
}

func (a AuthInternalService) getVersionFrSecret(secret string) (uint64, error) {
	if len(secret) != LENGTH_SECRET_STRING {
		return 0, errors.New(fmt.Sprintf("length secret must is %d", LENGTH_SECRET_STRING))
	}

	v, err := strconv.ParseUint(secret[20:40], 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("version secret %v is not uint, wrong format", secret[20:40]))
	}

	return v, nil
}

func (a AuthInternalService) checkAllConfigEnvBeforeRun() error {
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

	if os.Getenv(IS_USE_AUTH_INTERNAL_KEY) == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", IS_USE_AUTH_INTERNAL_KEY))
	}

	return nil
}

func (a *AuthInternalService) getSecretFrRemote() (AllSecret, error) {
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

func (a *AuthInternalService) UpdateSecretInLocal() error {
	s, errGetFrRemote := a.getSecretFrRemote()
	if errGetFrRemote != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetFrRemote.Error(),
		}).Error("getSecretFrRemote have error")
		return errGetFrRemote
	}

	fmt.Println("versionurrent %v \n", s.TimeUpdateSecret, s.SecretNearestOld, s.SecretCurrent)

	versionCurrent, errVC := a.getVersionFrSecret(s.SecretCurrent)
	if errVC != nil {
		log.WithFields(log.Fields{
			"errorMessage": errVC.Error(),
		}).Error("format secret have error")
		return errVC
	}

	versionNearestOld, errVNO := a.getVersionFrSecret(s.SecretNearestOld)
	if errVNO != nil {
		log.WithFields(log.Fields{
			"errorMessage": errVNO.Error(),
		}).Error("format secret have error")
		return errVNO
	}

	timeU, errParserTime := strconv.ParseUint(s.TimeUpdateSecret, 10, 64)
	if errParserTime != nil {
		log.WithFields(log.Fields{
			"errorMessage": errParserTime.Error(),
		}).Error("format secret have error")
		return errParserTime
	}

	//udpate
	a.secretCurrent = s.SecretCurrent
	a.secretNearestOld = s.SecretNearestOld
	a.versionCurrent = versionCurrent
	a.versionNearestOld = versionNearestOld
	a.timeUpdateSecret = timeU

	return nil
}

func (a *AuthInternalService) autoUpdateSecretInLocal() {
	for {
		err := a.UpdateSecretInLocal()
		if err != nil {
			log.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("autoUpdateSecretInLocal error")
			fmt.Println("autoUpdateSecretInLocal error")
		} else {
			log.Info("autoUpdateSecretInLocal success")
			fmt.Println("autoUpdateSecretInLocal success")
		}

		// Sleep() is the best for optimal cpu resource
		time.Sleep(time.Duration(a.timeUpdateSecret) * time.Millisecond)
	}
}

func (a *AuthInternalService) InitAuth() error {
	errEBR := a.checkAllConfigEnvBeforeRun()
	if errEBR != nil {
		log.WithFields(log.Fields{
			"errorMessage": errEBR.Error(),
		}).Error("missing config env before run packet %s", reflect.TypeOf(AuthInternalService{}).PkgPath())
		os.Exit(0)
	}

	err := a.UpdateSecretInLocal()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Panicf("init: update secret error")
	} else {
		log.Info("init: update secret in local success ")
	}

	authInternalMode := os.Getenv(IS_USE_AUTH_INTERNAL_KEY)
	if authInternalMode == IS_USE_AUTH_INTERNAL_VALUE {
		go func() {
			a.autoUpdateSecretInLocal()
		}()
	}

	return nil
}

func createFirstSecret() string {
	var version uint64
	a := AuthInternalService{}
	versionString := a.standardizedVersion(version)
	bodyString := uniuri.NewLen(LENGTH_BODY_STRING)

	secretNew := bodyString + versionString
	return secretNew
}

func (a AuthInternalService) standardizedVersion(v uint64) string {
	vNew := strconv.FormatUint(v, 10)
	nVewLen := len(vNew)

	if nVewLen < LENGTH_VERSION_STRING {
		for i := 0; i < LENGTH_VERSION_STRING-nVewLen; i++ {
			vNew = "0" + vNew
		}
	}

	return vNew
}

func (a AuthInternalService) ListEnvRequireSetupBeforeRunPacket() []string {
	return []string{
		SECRET_NAME_KEY,
		REGION_KEY,
		VERSION_STATE_KEY,
		IS_USE_AUTH_INTERNAL_KEY,
	}
}

func NewAuthInternalService() AuthInternalServiceInterface {
	return &AuthInternalService{}
}

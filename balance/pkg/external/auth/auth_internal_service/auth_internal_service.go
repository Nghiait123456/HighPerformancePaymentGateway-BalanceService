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

/**
  auth internal in micro service: high performance + security + handle all race conditions + zero down time when update new secret
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
		SecretNearestOld string `json:"SECRET_CURRENT"`
		SecretCurrent    string `json:"SECRET_NEAREST_OLD"`
		TimeUpdateSecret string `json:"TIME_UPDATE_SECRET"`
	}

	AuthInternalServiceInterface interface {
		/*
		  list fc for auth
		*/
		Auth(secretCheck string) (bool, AuthResponse) // auth to incoming request
		Init() error
		autoUpdateSecret()
		UpdateSecret() error
		getSecretFrRemote() (AllSecret, error)

		/*
		  list fc for auto change secret
		  just only one worker run this function
		  we have one worker run change, never run with many worker ==> race conditions
		*/
		ReplaceNewSecret()
		createNewSecret(currentSecret string) (string, bool)
		validateSecret(secret string) (bool, error)
	}
)

const (
	LENGTH_SECRET_STRING  = 40
	LENGTH_VERSION_STRING = 20
	LENGTH_BODY_STRING    = 20
	SECRET_NAME           = "AUTH_INTERNAL_SECRET_NAME"
	REGION                = "AUTH_INTERNAL_REGION"
	VERSION_STATTE        = ""

	SECRET_KEY_CURRENT = "SECRET_CURRENT"
	SECRET_NEAREST_OLD = "SECRET_NEAREST_OLD"
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
		err := a.UpdateSecret()
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
		return 0, errors.New("secret wrong format")
	}

	v, err := strconv.ParseUint(secret[20:40], 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("version secret %s is not uint, wrong format", secret[20:40]))
	}

	return v, nil

}

func (a *AuthInternalService) getSecretFrRemote() (AllSecret, error) {
	allS := AllSecret{}
	sAws := env.NewAwsManagerSecret()
	secretName := os.Getenv(SECRET_NAME)
	if secretName == "" {
		log.WithFields(log.Fields{
			"secretName": "",
		}).Error("Secret name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	region := os.Getenv(REGION)
	if region == "" {
		log.WithFields(log.Fields{
			region: "",
		}).Error("region name is empty")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	versionState := ""
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

	return allS, nil
}

func (a *AuthInternalService) UpdateSecret() error {
	s, errGetFrRemote := a.getSecretFrRemote()
	if errGetFrRemote != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetFrRemote.Error(),
		}).Panicf("getSecretFrRemote have error")
		panic(fmt.Sprintf("getSecretFrRemote have error: %s", errGetFrRemote.Error()))
	}

	versionCurrent, errVC := a.getVersionFrSecret(s.SecretCurrent)
	if errVC != nil {
		log.WithFields(log.Fields{
			"errorMessage": errVC.Error(),
		}).Panicf("format secret have error")
		panic(fmt.Sprintf("format secret_current from auth_internal service have error: %s", errVC.Error()))
	}

	versionNearestOld, errVNO := a.getVersionFrSecret(a.secretNearestOld)
	if errVNO != nil {
		log.WithFields(log.Fields{
			"errorMessage": errVNO.Error(),
		}).Panicf("format secret have error")
		panic(fmt.Sprintf("format secretNearestOld from auth_internal service have error %s", errVNO.Error()))
	}

	timeU, errParserTime := strconv.ParseUint(s.TimeUpdateSecret, 10, 64)
	if errParserTime != nil {
		log.WithFields(log.Fields{
			"errorMessage": errParserTime.Error(),
		}).Panicf("format secret have error")
		panic(fmt.Sprintf("format secret TimeUpdateSecret from auth_internal service have error %s", errParserTime.Error()))
	}

	//udpate
	a.secretCurrent = s.SecretCurrent
	a.secretNearestOld = s.SecretNearestOld
	a.versionCurrent = versionCurrent
	a.versionNearestOld = versionNearestOld
	a.timeUpdateSecret = timeU

	return nil
}

func (a *AuthInternalService) autoUpdateSecret() {
	for {
		err := a.UpdateSecret()
		if err != nil {
			log.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Panicf("update secret error")
		}

		// Sleep() is the best for optimal cpu resource
		time.Sleep(time.Duration(a.timeUpdateSecret) * time.Millisecond)
	}
}

func (a *AuthInternalService) Init() error {
	err := a.UpdateSecret()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Panicf("update secret error")
		panic(fmt.Sprintf("update secret error: %s", err.Error()))
	}

	go func() {
		a.autoUpdateSecret()
	}()

	return nil
}

func (a AuthInternalService) createNewSecret(currentSecret string) (string, error) {
	version, errV := a.getVersionFrSecret(currentSecret)
	if errV != nil {
		log.WithFields(log.Fields{
			"errorMessage": errV.Error(),
		}).Panicf("format secret have error")
		panic(fmt.Sprintf("format secret_current from auth_internal service have error: %s", errV.Error()))
	}

	version++
	versionString := a.standardizedVersion(version)
	bodyString := uniuri.NewLen(LENGTH_BODY_STRING)

	secretNew := bodyString + versionString
	return secretNew, nil
}

func (a AuthInternalService) standardizedVersion(v uint64) string {
	vNew := strconv.FormatUint(v, 10)
	if len(vNew) < LENGTH_VERSION_STRING {
		for i := 0; i < LENGTH_VERSION_STRING-len(vNew); i++ {
			vNew = "0" + vNew
		}
	}

	return vNew
}

func (a *AuthInternalService) ReplaceNewSecret() {
	s, errGetFrRemote := a.getSecretFrRemote()
	if errGetFrRemote != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetFrRemote.Error(),
		}).Panicf("getSecretFrRemote have error")
		panic(fmt.Sprintf("getSecretFrRemote have error: %s", errGetFrRemote.Error()))
	}

	newSecret, errNS := a.createNewSecret(s.SecretCurrent)
	if errNS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errNS.Error(),
		}).Panicf("createNewSecret have error")
		panic(fmt.Sprintf("createNewSecret have error: %s", errNS.Error()))
	}

	secretUpdate := AllSecret{
		SecretCurrent:    newSecret,
		SecretNearestOld: s.SecretCurrent,
		TimeUpdateSecret: s.TimeUpdateSecret,
	}

}

func NewAuthInternalService() AuthInternalServiceInterface {
	return &AuthInternalService{}
}

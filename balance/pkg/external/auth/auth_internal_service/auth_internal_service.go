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
		env               env.AuthInternalServiceConfigInterface
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
		Auth(secretCheck string) (bool, AuthResponse) // auth to incoming request
		Init() error
		autoUpdateSecretInLocal()
		updateSecretInLocal() error
		TryUpdateSecretInLocal(numberTry uint) error
		getSecretFrRemote() (AllSecret, error)
		checkAllConfigEnvBeforeRun() error
	}
)

const (
	LENGTH_SECRET_STRING              = 40
	LENGTH_VERSION_STRING             = 20
	LENGTH_BODY_STRING                = 20
	TIME_UPDATE_SECRET_DEFAULT        = 10000 //ms
	NUMBER_TRY_UPDATE_SECRET_IN_LOCAL = 20

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
		err := a.TryUpdateSecretInLocal(NUMBER_TRY_UPDATE_SECRET_IN_LOCAL)
		if err != nil {
			log.WithFields(log.Fields{
				"errMessage": err.Error(),
			}).Error("TryUpdateSecretInLocal secret error")
			panic(fmt.Sprintf("TryUpdateSecretInLocal: %s", err.Error()))
			os.Exit(0)
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
	if a.env.SecretName() == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", env.AUTH_INTERNER_SERVICE_SECRET_NAME_KEY))
	}

	if a.env.Region() == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", env.AUTH_INTERNER_SERVICE_REGION_KEY))
	}

	if a.env.VersionState() == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", env.AUTH_INTERNER_SERVICE_VERSION_STATE_KEY))
	}

	if a.env.IsUseAuthInternal() == "" {
		return errors.New(fmt.Sprintf("missing config env key %s", env.AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL_KEY))
	}

	return nil
}

func (a *AuthInternalService) getSecretFrRemote() (AllSecret, error) {
	allS := AllSecret{}
	sAws := env.NewAwsManagerSecret()
	sAws.Init(a.env.SecretName(), a.env.Region(), a.env.VersionState())

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

func (a *AuthInternalService) TryUpdateSecretInLocal(numberTry uint) error {
	success := false
	var errGL error

	for i := uint(0); i < numberTry; i++ {
		err := a.updateSecretInLocal()
		errGL = err
		if err != nil {
			log.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("dont get secret from remote")
		} else {
			success = true
			break
		}
	}

	if success != true {
		return errGL
	}

	return nil
}

func (a *AuthInternalService) updateSecretInLocal() error {
	s, errGetFrRemote := a.getSecretFrRemote()
	if errGetFrRemote != nil {
		log.WithFields(log.Fields{
			"errorMessage": errGetFrRemote.Error(),
		}).Error("getSecretFrRemote have error")
		return errGetFrRemote
	}

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
		fmt.Println("autoUpdateSecretInLocal...")
		err := a.TryUpdateSecretInLocal(NUMBER_TRY_UPDATE_SECRET_IN_LOCAL)
		if err != nil {
			log.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("TryUpdateSecretInLocal error")
			panic(fmt.Sprintf("TryUpdateSecretInLocal error: %s", err.Error()))
			os.Exit(0)
		} else {
			log.Info("autoUpdateSecretInLocal success")
			fmt.Println("autoUpdateSecretInLocal success")
		}

		// Sleep() is the best for optimal cpu resource
		time.Sleep(time.Duration(a.timeUpdateSecret) * time.Millisecond)
	}
}

func (a *AuthInternalService) Init() error {
	errEBR := a.checkAllConfigEnvBeforeRun()
	if errEBR != nil {
		log.WithFields(log.Fields{
			"errorMessage": errEBR.Error(),
		}).Errorf("missing config env before run packet %s", reflect.TypeOf(AuthInternalService{}).PkgPath())
		panic(fmt.Sprintf("missing config env before run packet %s", reflect.TypeOf(AuthInternalService{}).PkgPath()))
		os.Exit(0)
	}

	errTUS := a.TryUpdateSecretInLocal(NUMBER_TRY_UPDATE_SECRET_IN_LOCAL)
	if errTUS != nil {
		log.WithFields(log.Fields{
			"errorMessage": errTUS.Error(),
		}).Error("init: TryUpdateSecretInLocal secret error")
		panic(fmt.Sprintf("TryUpdateSecretInLocal error : %s", errTUS.Error()))
	} else {
		log.Info("init: update secret in local success ")
	}

	if a.env.IsUseAuthInternal() == IS_USE_AUTH_INTERNAL_VALUE {
		fmt.Println("aaa")
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

func NewAuthInternalService(env env.AuthInternalServiceConfigInterface) AuthInternalServiceInterface {
	return &AuthInternalService{
		env: env,
	}
}

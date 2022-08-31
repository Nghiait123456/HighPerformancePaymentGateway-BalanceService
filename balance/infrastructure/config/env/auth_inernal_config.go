package env

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	AuthInternalServiceConfig struct {
		secretName            string
		region                string
		versionState          string
		isUseAuthInternal     string
		isUseAutoChangeSecret string
	}

	AuthInternalServiceConfigInterface interface {
		BaseConfigInterface
		SecretName() string
		Region() string
		VersionState() string
		IsUseAuthInternal() string
		IsUseAutoChangeSecret() string
	}
)

const (
	AUTH_INTERNER_SERVICE_SECRET_NAME_KEY                  = "AUTH_INTERNER_SERVICE_SECRET_NAME"
	AUTH_INTERNER_SERVICE_REGION_KEY                       = "AUTH_INTERNER_SERVICE_REGION"
	AUTH_INTERNER_SERVICE_VERSION_STATE_KEY                = "AUTH_INTERNER_SERVICE_VERSION_STATE"
	AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL_KEY         = "AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL"
	AUTH_INTERNAL_SERVICE_IS_AUTO_CHANGE_SECRET_REMOTE_KEY = "AUTH_INTERNAL_SERVICE_IS_AUTO_CHANGE_SECRET_REMOTE"
)

func (a *AuthInternalServiceConfig) Load() {
	a.secretName = os.Getenv(AUTH_INTERNER_SERVICE_SECRET_NAME_KEY)
	if a.secretName == "" {
		log.WithFields(log.Fields{
			"secretName": "",
		}).Error("Secret name dont have config in env")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	a.region = os.Getenv(AUTH_INTERNER_SERVICE_REGION_KEY)
	if a.region == "" {
		log.WithFields(log.Fields{
			"region": "",
		}).Error("region name dont have config in env")
		panic("Secret name of server auth internal is empty")
		os.Exit(0)
	}

	versionState, errVST := os.LookupEnv(AUTH_INTERNER_SERVICE_VERSION_STATE_KEY)
	if errVST != true {
		log.WithFields(log.Fields{
			"versionState": "",
		}).Error("versionState dont have config in env")
		panic("versionState dont have config in env")
		os.Exit(0)
	}
	a.versionState = versionState

	a.isUseAuthInternal = os.Getenv(AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL_KEY)
	if a.isUseAuthInternal == "" {
		log.WithFields(log.Fields{
			"isUseAuthInternal": "",
		}).Error("AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL_KEY dont have config in env")
		panic("AUTH_INTERNER_SERVICE_IS_USE_AUTH_INTERNAL_KEY dont have config in env")
		os.Exit(0)
	}

	a.isUseAutoChangeSecret = os.Getenv(AUTH_INTERNAL_SERVICE_IS_AUTO_CHANGE_SECRET_REMOTE_KEY)
	if a.isUseAutoChangeSecret == "" {
		log.WithFields(log.Fields{
			"isUseAuthInternal": "",
		}).Error("AUTH_INTERNAL_SERVICE_IS_AUTO_CHANGE_SECRET_REMOTE_KEY dont have config in env")
		panic("AUTH_INTERNAL_SERVICE_IS_AUTO_CHANGE_SECRET_REMOTE_KEY dont have config in env")
		os.Exit(0)
	}

}

func (a AuthInternalServiceConfig) SecretName() string {
	return a.secretName
}

func (a AuthInternalServiceConfig) Region() string {
	return a.region
}
func (a AuthInternalServiceConfig) VersionState() string {
	return a.versionState
}

func (a AuthInternalServiceConfig) IsUseAuthInternal() string {
	return a.isUseAuthInternal
}

func (a AuthInternalServiceConfig) IsUseAutoChangeSecret() string {
	return a.isUseAutoChangeSecret
}

func NewAuthInternalServiceConfig() AuthInternalServiceConfigInterface {
	return &AuthInternalServiceConfig{}
}

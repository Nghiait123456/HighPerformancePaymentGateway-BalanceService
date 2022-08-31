package env

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"os"
)

type (
	GlobalConfig struct {
		authInternalServiceConfig AuthInternalServiceConfigInterface
	}

	ConfigEnv struct {
		FilePatchLogInLocal string
	}

	GlobalConfigInterface interface {
		loadEnvLocal(c ConfigEnv)
		loadEnvDev()
		loadEnvStaging()
		loadEnvProduct()
		loadEnv(c ConfigEnv)
		mapEnvToStruct()
		validateEnvironmentType() (string, error)
		Init(c ConfigEnv)
		ListEnvRequireSetUpBeforeInit() []string
		CheckEnvRequireSetUpBeforeInit()

		/**
		  function for get
		*/
		AuthInternalServiceConfig() AuthInternalServiceConfigInterface
	}
)

const (
	ENV_LOCAL       = "local"
	ENV_DEV         = "dev"
	ENV_STAGING     = "staging"
	ENV_PRODUCT     = "product"
	ENV_ENVIRONMENT = "ENVIRONMENT_BALANCE_SERVICE"

	AWS_SECRET_NAME_GLOBAL_KEY   = "AWS_SECRET_NAME_GLOBAL_KEY"
	AWS_REGION_GLOBAL_KEY        = "AWS_REGION_GLOBAL_KEY"
	AWS_VERSION_STATE_GLOBAL_KEY = "AWS_VERSION_STATE_GLOBAL_KEY"
)

func (g GlobalConfig) AllEnvironment() []string {
	return []string{ENV_LOCAL, ENV_DEV, ENV_STAGING, ENV_PRODUCT}
}

func (g GlobalConfig) IsEnvironmentValid(v string) bool {
	return slices.Contains(g.AllEnvironment(), v)
}

func (g *GlobalConfig) validateEnvironmentType() (string, error) {
	evm := os.Getenv(ENV_ENVIRONMENT)
	if len(evm) == 0 {
		errM := fmt.Sprintf("ENV_ENVIRONMENT with key %s empty or not exits", ENV_ENVIRONMENT)
		panic(errM)
		os.Exit(0)
	}

	if !g.IsEnvironmentValid(evm) {
		errM := fmt.Sprintf("ENV_ENVIRONMENT with key %s with value %s not valid", ENV_ENVIRONMENT, evm)
		panic(errM)
		os.Exit(0)
	}

	return evm, nil
}

func (g *GlobalConfig) loadEnv(c ConfigEnv) {
	evm, err := g.validateEnvironmentType()
	if err != nil {
		panic(err.Error())
		os.Exit(0)
	}

	switch evm {
	case ENV_LOCAL:
		g.loadEnvLocal(c)
		return
	case ENV_DEV:
		g.loadEnvDev()
		return
	case ENV_PRODUCT:
		g.loadEnvProduct()
		return
	}
}

func (g *GlobalConfig) loadEnvLocal(c ConfigEnv) {
	//load env
	errL := godotenv.Load(c.FilePatchLogInLocal)
	if errL != nil {
		log.WithFields(log.Fields{
			"errMessage": errL.Error(),
		}).Error("load env file error")
		panic(fmt.Sprintf("load env file error %s", errL.Error()))
		os.Exit(0)
	}
}

func (g *GlobalConfig) loadEnvDev() {
	// never save secret to file, always get from api other service
	aws := NewAwsManagerSecret()

	// init
	secretName := os.Getenv(AWS_SECRET_NAME_GLOBAL_KEY)
	if secretName == "" {
		panic("don't exits AWS_SECRET_NAME_GLOBAL_KEY in dev")
		os.Exit(0)
	}

	region := os.Getenv(AWS_REGION_GLOBAL_KEY)
	if region == "" {
		panic("don't exits AWS_REGION_GLOBAL_KEY in dev")
		os.Exit(0)
	}

	versionState, errV := os.LookupEnv(AWS_VERSION_STATE_GLOBAL_KEY)
	if errV != true {
		panic("don't exits AWS_REGION_GLOBAL_KEY in dev")
		os.Exit(0)
	}

	aws.Init(secretName, region, versionState)

	secretS, errGS := aws.GetSecret()
	if errGS != nil {
		panic(fmt.Sprintf("get secret from aws manager secret in dev error: %s", errGS.Error()))
		os.Exit(0)
	}

	serretM := map[string]string{}
	errCVS := json.Unmarshal([]byte(secretS), &serretM)
	if errCVS != nil {
		panic(fmt.Sprintf("secret wrong format when convert form json to map error: %s", errCVS.Error()))
		os.Exit(0)
	}

	// merger
	for k, v := range serretM {
		err := os.Setenv(k, v)
		if err != nil {
			panic(fmt.Sprintf("set env have error: %s", err.Error()))
			os.Exit(0)
		}
	}
}

func (g *GlobalConfig) loadEnvStaging() {
	//return nil
}

func (g *GlobalConfig) loadEnvProduct() {
	//return nil
}

func (g *GlobalConfig) contructAllChildStruct() {
	g.authInternalServiceConfig = NewAuthInternalServiceConfig()
}
func (g *GlobalConfig) mapEnvToStruct() {
	g.authInternalServiceConfig.Load()
}

func (g *GlobalConfig) Init(c ConfigEnv) {
	g.CheckEnvRequireSetUpBeforeInit()
	g.loadEnv(c)
	g.contructAllChildStruct()
	g.mapEnvToStruct()
}

func (g GlobalConfig) ListEnvRequireSetUpBeforeInit() []string {
	return []string{
		ENV_ENVIRONMENT,
		AWS_SECRET_NAME_GLOBAL_KEY,
		AWS_REGION_GLOBAL_KEY,
		AWS_VERSION_STATE_GLOBAL_KEY,
	}
}

func (g GlobalConfig) CheckEnvRequireSetUpBeforeInit() {
	envEnviro := os.Getenv(ENV_ENVIRONMENT)
	if envEnviro == "" {
		panic("don't exits ENV_ENVIRONMENT")
		os.Exit(0)
	}

	if envEnviro != ENV_LOCAL {
		secretName := os.Getenv(AWS_SECRET_NAME_GLOBAL_KEY)
		if secretName == "" {
			panic("don't exits AWS_SECRET_NAME_GLOBAL_KEY")
			os.Exit(0)
		}

		region := os.Getenv(AWS_REGION_GLOBAL_KEY)
		if region == "" {
			panic("don't exits AWS_REGION_GLOBAL_KEY")
			os.Exit(0)
		}

		_, errV := os.LookupEnv(AWS_VERSION_STATE_GLOBAL_KEY)
		if errV != true {
			panic("don't exits AWS_REGION_GLOBAL_KEY")
			os.Exit(0)
		}
	}
}

func (g GlobalConfig) AuthInternalServiceConfig() AuthInternalServiceConfigInterface {
	return g.authInternalServiceConfig
}
func NewGlobalConfig() GlobalConfigInterface {
	g := GlobalConfig{}
	return &g
}

package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
	"log"
	"os"
)

type (
	GlobalConfig struct {
		DefaultDBSqlUserName string `env:"DEFAULT_DB_SQL_USERNAME,notEmpty"`
		DefaultDBSqlPassWord string `env:"DEFAULT_DB_SQL_PASSWORD,notEmpty"`
		DefaultDBSqlLink     string `env:"DEFAULT_DB_SQL_LINK,notEmpty"`
		DefaultDBSqlNameDB   string `env:"DEFAULT_DB_SQL_NAME_DB,notEmpty"`
	}

	DBConfigSqlDefault struct {
		UserName string
		PassWord string
		Link     string
		NameDB   string
	}

	DBConfigSqlDefaultInterface interface {
		GetDSN() (string, error)
	}

	GlobalConfigInterface interface {
		LoadConfig() error
		MapConfigToStruct() error
		DetectEnvironment() string
		AllEnvironment() []string
	}
)

const (
	ENV_LOCAL       = "local"
	ENV_DEV         = "dev"
	ENV_STAGING     = "staging"
	ENV_PRODUCT     = "product"
	ENV_ENVIRONMENT = "ENVIRONMENT_BALANCE_SERVICE"

	AWS_SECRET_NAME   = "AWS_SECRET_NAME"
	AWS_REGION        = "AWS_REGION"
	AWS_VERSION_STATE = "AWS_VERSION_STATE"
)

func (g GlobalConfig) AllEnvironment() []string {
	return []string{ENV_LOCAL, ENV_DEV, ENV_STAGING, ENV_PRODUCT}
}

func (g GlobalConfig) IsEnvironmentValid(v string) bool {
	return slices.Contains(g.AllEnvironment(), v)
}

func (g *GlobalConfig) DetectEnvironment() (string, error) {
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

func (g *GlobalConfig) LoadConfig() error {
	evm, err := g.DetectEnvironment()
	if err != nil {
		panic(err.Error())
		os.Exit(0)
	}

	switch evm {
	case ENV_LOCAL:
		g.LoadEnvLocal()

	}

	return nil
}

func (g *GlobalConfig) LoadEnvLocal() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return nil
}

func (g *GlobalConfig) LoadEnvDev() error {
	// never save secret to file, always get from api other service
	aws := NewAwsManagerSecret()

	// init
	secretName := os.Getenv(AWS_SECRET_NAME)
	if secretName == "" {
		panic("don't exits AWS_SECRET_NAME in dev")
		os.Exit(0)
	}

	region := os.Getenv(AWS_REGION)
	if region == "" {
		panic("don't exits AWS_REGION in dev")
		os.Exit(0)
	}

	versionState := os.Getenv(AWS_VERSION_STATE)
	if versionState == "" {
		panic("don't exits AWS_REGION in dev")
		os.Exit(0)
	}

	aws.Init(secretName, region, versionState)

	//get
	//secretV, err := aws.GetSecret()
	//if err != nil {
	//	panic(fmt.Sprintf("get secret from aws manager secret in dev error, %s", err.Error()))
	//	os.Exit(0)
	//}

	// todo mapping secret to struct

	return nil

}

func (g *GlobalConfig) LoadEnvStaging() error {
	return nil
}

func (g *GlobalConfig) LoadEnvProduct() error {
	return nil
}

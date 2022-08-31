package env

//package main
//
//import (
//	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
//	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
//	"os"
//	"time"
//)

//func main() {
//	log_init.Init(log_init.Log{
//		TypeFormat: log_init.TYPE_FORMAT_TEXT,
//		TypeOutput: log_init.TYPE_OUTPUT_FILE,
//		LinkFile:   "balance/infrastructure/log/log_file/log.log",
//	})
//
//	gCf := TestEnvGlobal()
//	TestAutoChangeSecret(gCf.AuthInternalServiceConfig())
//
//	time.Sleep(1000 * time.Second)
//}

//func TestEnvGlobal() env.GlobalConfigInterface {
//	os.Setenv(env.ENV_ENVIRONMENT, env.ENV_DEV)
//	os.Setenv(env.AWS_SECRET_NAME_GLOBAL_KEY, "payment-balance-service-qwedjfndasndajndn12")
//	os.Setenv(env.AWS_REGION_GLOBAL_KEY, "ap-southeast-1")
//	os.Setenv(env.AWS_VERSION_STATE_GLOBAL_KEY, "AWSCURRENT")
//
//	gCf := env.NewGlobalConfig()
//	gCf.Init(env.ConfigEnv{
//		FilePatchLogInLocal: "balance/infrastructure/config/env/.env",
//	})
//
//	return gCf
//}

//func TestAutoChangeSecret(ai env.AuthInternalServiceConfigInterface) {
//firstInit
//err := auth_internal_service.FirstTimeContructSecret()
//if err != nil {
//	fmt.Println("contrcut secret firsttime erro : %s", err.Error())
//} else {
//	fmt.Println("contruct secret in firsttime success")
//}

//auto change secret
//auto := auth_internal_service.NewAutoChangeSecret(ai)
//err := auto.Init()
//if err != nil {
//	fmt.Println("init error : %s", err.Error())
//} else {
//	fmt.Println("init success")
//}

//test auth
//auth := auth_internal_service.NewAuthInternalService(ai)
//auth.Init()
//fmt.Println(auth.Auth("42zPQ0sIAmjtHSTRANmc00000000000000000004"))

// test integrated

//auto change secret
//go func() {
//	auto := auth_internal_service.NewAutoChangeSecret(ai)
//	err := auto.Init()
//	if err != nil {
//		fmt.Println("init error : %s", err.Error())
//	} else {
//		fmt.Println("init success")
//	}
//}()
//
////init auth
//auth := auth_internal_service.NewAuthInternalService(ai)
//auth.Init()
//go func() {
//	for {
//		fmt.Println(auth.Auth("K38xvC2XthUC6fHSxz6W00000000000000000006"))
//		time.Sleep(5 * time.Second)
//	}
//}()
//}

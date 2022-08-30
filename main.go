// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html
package main

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/auth/auth_internal_service"
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
	"os"
	"time"
)

func main() {
	log_init.Init(log_init.Log{
		TypeFormat: log_init.TYPE_FORMAT_TEXT,
		TypeOutput: log_init.TYPE_OUTPUT_FILE,
		LinkFile:   "balance/infrastructure/log/log_file/log.log",
	})

	TestAuthServiceInternal()

	time.Sleep(1000 * time.Second)

}

func TestAuthServiceInternal() {
	os.Setenv(auth_internal_service.SECRET_NAME_KEY, "payment-balance-service-qwedjfndasndajndn12")
	os.Setenv(auth_internal_service.REGION_KEY, "ap-southeast-1")
	os.Setenv(auth_internal_service.VERSION_STATE_KEY, "")
	os.Setenv(auth_internal_service.IS_USE_AUTH_INTERNAL_KEY, auth_internal_service.IS_USE_AUTH_INTERNAL_VALUE)
	os.Setenv(auth_internal_service.IS_AUTO_CHANGE_SECRET_REMOTE_KEY, auth_internal_service.IS_AUTO_CHANGE_SECRET_REMOTE_VALUE)

	//auto := auth_internal_service.NewAutoChangeSecret()
	//auto.Init()

	//first time contruct
	//rs := auth_internal_service.FirstTimeContructSecret()
	//if rs != nil {
	//	panic("construct secret in first time error")
	//	os.Exit(0)
	//}
	//fmt.Printf("init first time contruct success")

	//authIS := auth_internal_service.NewAuthInternalService()
	//authIS.InitAuth()
	//rs, err := authIS.Auth("x8kx3QFk9uxwMBFFT1zS00000000000000000008")
	//fmt.Printf("rs = %v, err = %v", rs, err)

}

func UpdateAndGetSecret() {
	secretName := "payment-balance-service-qwedjfndasndajndn12"
	region := "ap-southeast-1"
	versionState := ""

	update := `{"TEST1":"ssssssssssss3333333333333"}`

	awsS := env.NewAwsManagerSecret()
	awsS.Init(secretName, region, versionState)
	errU := awsS.UpdateSecretString(update)
	if errU != nil {
		if awsS.IsErrorCustomOfAws(errU) {
			errN := awsS.ConvertToErrorAws(errU)
			fmt.Printf("errorCode %s errMessage %s ", errN.Code(), errN.Message())
		} else {
			fmt.Printf(" errMessage %s ", errU.Error())
		}
	} else {
		fmt.Printf("update success \n")
	}

	value, errG := awsS.GetSecret()
	if errG != nil {
		if awsS.IsErrorCustomOfAws(errG) {
			errN := awsS.ConvertToErrorAws(errG)
			fmt.Printf("errorCode %s errMessage %s ", errN.Code(), errN.Message())
		} else {
			fmt.Printf(" errMessage %s ", errG.Error())
		}
	}

	fmt.Printf("get success:  data get %s", value)

}

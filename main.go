// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html
package main

import (
	"fmt"
	"github.com/high-performance-payment-gateway/balance-service/balance/infrastructure/config/env"
)

func main() {
	UpdateAndGetSecret()
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

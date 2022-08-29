package env

import "fmt"

/*
 example update and get secret
*/
func updateAndGetSecret() {
	secretName := "payment-balance-service-qwedjfndasndajndn12"
	region := "ap-southeast-1"
	versionState := ""

	update := `{"TEST1":"ssssssssssss3333333333333"}`

	awsS := NewAwsManagerSecret()
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

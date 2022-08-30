package auth_internal_service

import "os"

func TestAuthServiceInternal() {
	os.Setenv(SECRET_NAME_KEY, "payment-balance-service-qwedjfndasndajndn12")
	os.Setenv(REGION_KEY, "ap-southeast-1")
	os.Setenv(VERSION_STATE_KEY, "")
	os.Setenv(IS_USE_AUTH_INTERNAL_KEY, IS_USE_AUTH_INTERNAL_VALUE)
	os.Setenv(IS_AUTO_CHANGE_SECRET_REMOTE_KEY, IS_AUTO_CHANGE_SECRET_REMOTE_VALUE)

	// test auto change secret
	//auto := auth_internal_service.NewAutoChangeSecret()
	//auto.Init()

	//test first time contruct
	//rs := auth_internal_service.FirstTimeContructSecret()
	//if rs != nil {
	//	panic("construct secret in first time error")
	//	os.Exit(0)
	//}
	//fmt.Printf("init first time contruct success")

	// test auth
	//authIS := auth_internal_service.NewAuthInternalService()
	//authIS.InitAuth()
	//rs, err := authIS.Auth("x8kx3QFk9uxwMBFFT1zS00000000000000000008")
	//fmt.Printf("rs = %v, err = %v", rs, err)

}

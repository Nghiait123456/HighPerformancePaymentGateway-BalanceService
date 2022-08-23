package error_base

import (
	"fmt"
)

func handle() {
	err := testError()
	if err != nil {
		if IsErrorBase(err) {
			errB := GetErrorBase(err)
			//todo handle Error Code Custom
			//switch errB.Code() {
			//case :
			//}
			fmt.Println(errB.GetCode())
		} else {
			// todo handle defaut error golang
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("error = nil")
	}
}

func testError() error {
	return New(3, "test")
}

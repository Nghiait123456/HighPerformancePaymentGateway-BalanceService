package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidateFormExample = ValidateBase

func Example() {
	//v := ValidateFormExample{}
	//v.SetMapRuleToMessage(make(MapMessage))
	//v.SetValidate(validator.New())

	v := NewBaseValidate()

	v.ResignValidateCustom("minLength", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	v.ResignValidateCustom("maxLength", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) < 25
	})

	message := make(map[string]string)
	message["required"] = "filed is required"
	message["minLength"] = "field smaller minLength"
	message["maxLength"] = "field Greater MaxLength"
	message["gte"] = "field over gte"
	message["lte"] = "field over lte"

	v.SetMessageForRule(message)

	type User struct {
		FirstName string `json:"fname" validate:"required,minLength"`
		LastName  string `json:"lname" validate:"required,maxLength"`
		Age       uint8  `validate:"gte=20,lte=65"`
	}

	user := &User{
		FirstName: "A",
		LastName:  "Te8888888888888888888888888888888888888888888888888st",
		Age:       75,
	}

	err := v.Validate().Struct(user)

	if err == nil {
		fmt.Println("dont have error")
	}

	messageErr, errCv := v.ConvertErrorValidate(err)
	if errCv != nil {
		fmt.Println("ConvertErrorValidate error")
	} else {
		errS, Message := v.ShowErrors(messageErr, DefaultShowErrors)
		if errS != nil {
			fmt.Println("ShowErrors error")
		} else {
			fmt.Println(fmt.Sprintf("ShowErrors success, error %v", Message))
		}
	}
}

/**
error {"ListError":{"Age":{"Field":"Age","Rule":"lte","Message":"field over lte","ParamRule":"65"},"FirstName":{"Field":"FirstName","Rule":"minLength","Message":"field smaller minLength","ParamRule":""},"LastName":{"Field":"LastName","Rule":"maxLength","Message":"field Greater MaxLength","ParamRule":""}}}
*/

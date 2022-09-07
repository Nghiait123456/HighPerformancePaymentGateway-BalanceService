package validate

import "github.com/go-playground/validator/v10"

type (
	ErrorBase struct {
		Field      string
		Rule       string
		Message    string
		ValueError interface{}
		ParamRule  string
		Raw        validator.FieldError
	}
)

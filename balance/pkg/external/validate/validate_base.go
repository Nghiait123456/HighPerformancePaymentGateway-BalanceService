package validate

import (
	"github.com/go-playground/validator/v10"
)

type (
	ValidateBase struct {
		validate         *validator.Validate
		mapRuleToMessage MapMessage // [rule]"message"
	}

	ShowError      = func(MessageErrors) (error, interface{})
	MessageErrors  map[string]ErrorBase //[field]ErrorBase
	ValidateCustom = validator.Func
	MapMessage     = map[string]string

	ValidateBaseInterface interface {
		ResignValidateCustom(ruleNameCustom string, vc ValidateCustom)
		ConvertErrorValidate(err error) (error, MessageErrors)
		ShowErrors(m MessageErrors, s ShowError) (error, interface{})
		SetMessageForRule(m map[string]string)
		Validate() *validator.Validate
		SetValidate(va *validator.Validate)
		SetMapRuleToMessage(vM MapMessage)
	}
)

func (v *ValidateBase) ConvertErrorValidate(err error) (error, MessageErrors) {
	MessageErrors := make(MessageErrors)

	if errIVE, ok := err.(*validator.InvalidValidationError); ok {
		return errIVE, MessageErrors
	}

	for _, errV := range err.(validator.ValidationErrors) {
		ruleCode := errV.ActualTag()
		var message string
		if va, ok := v.mapRuleToMessage[ruleCode]; ok {
			message = va
		}

		eb := ErrorBase{
			Field:      errV.StructField(),
			Rule:       errV.ActualTag(),
			Message:    message,
			ParamRule:  errV.Param(),
			ValueError: errV.Value(),
			Raw:        errV,
		}

		MessageErrors[errV.StructField()] = eb
	}

	return nil, MessageErrors
}

func (v *ValidateBase) ResignValidateCustom(ruleNameCustom string, vc ValidateCustom) {
	v.validate.RegisterValidation(ruleNameCustom, vc)
}

func (v *ValidateBase) ShowErrors(m MessageErrors, s ShowError) (error, interface{}) {
	return s(m)
}

func (v *ValidateBase) SetMessageForRule(m map[string]string) {
	v.mapRuleToMessage = m
}

func (v *ValidateBase) Validate() *validator.Validate {
	return v.validate
}

func (v *ValidateBase) SetValidate(va *validator.Validate) {
	v.validate = va
}

func (v *ValidateBase) SetMapRuleToMessage(vM MapMessage) {
	v.mapRuleToMessage = vM
}

func NewBaseValidate() ValidateBaseInterface {
	return &ValidateBase{}
}

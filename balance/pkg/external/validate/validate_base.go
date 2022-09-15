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
		ConvertErrorValidate(err error) (MessageErrors, error)
		ShowErrors(m MessageErrors, s ShowError) (error, interface{})
		SetMessageForRule(m MapMessage)
		Validate() *validator.Validate
		SetValidate(va *validator.Validate)
		SetMapRuleToMessage(vM MapMessage)
	}
)

func (v *ValidateBase) ConvertErrorValidate(err error) (MessageErrors, error) {
	mErrors := make(MessageErrors)

	if errIVE, ok := err.(*validator.InvalidValidationError); ok {
		return mErrors, errIVE
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

		mErrors[errV.StructField()] = eb
	}

	return mErrors, nil
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
	v := ValidateBase{}
	v.SetMapRuleToMessage(make(MapMessage))
	v.SetValidate(validator.New())

	return &v
}

package validate

import "encoding/json"

type (
	ListErrorsDefaultShow map[string]OneErrorDefaultShow
	DefaultShowError      struct {
		ListError ListErrorsDefaultShow
	}

	OneErrorDefaultShow struct {
		Field      string
		Rule       string
		Message    string
		ParamRule  string
		ValueError any
		Tag        string
	}
)

func DefaultShowErrors(mE MessageErrors) (error, interface{}) {
	ListES := ListErrorsDefaultShow{}

	for _, v := range mE {
		oneErr := OneErrorDefaultShow{
			Field:     v.Field,
			Rule:      v.Rule,
			Message:   v.Message,
			ParamRule: v.ParamRule,
		}

		ListES[v.Field] = oneErr
	}

	showE := DefaultShowError{
		ListError: ListES,
	}

	rs, errCV := json.Marshal(showE)
	if errCV != nil {
		return errCV, ""
	}

	return nil, string(rs)
}

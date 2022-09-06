package error_base

import "fmt"

type (
	// A BaseError wraps the Code and message which defines an error.
	ErrorBase struct {
		Code      int
		Message   string
		Signature string // use detect type error
	}

	//One Error Custom matching error default golang
	Error interface {
		// Satisfy the generic error interface.
		error
		GetCode() int
		GetMessage() string
		GetSignature() string
	}
)

func (b *ErrorBase) GetCode() int {
	return b.Code
}

func (b *ErrorBase) GetMessage() string {
	return b.Message
}

func (b *ErrorBase) GetSignature() string {
	return b.Signature
}

func newBaseError(code int, message string) *ErrorBase {
	b := &ErrorBase{
		Code:    code,
		Message: message,
	}

	return b
}

// Satisfies the error interface.
func (b *ErrorBase) Error() string {
	return SprintError(b.Code, b.Message)
}

func SprintError(code int, message string) string {
	msg := fmt.Sprintf("error Code: %d , error Message:  %s", code, message)
	return msg
}

func IsErrorBase(e error) bool {
	if _, ok := e.(Error); ok {
		return true
	}

	return false
}

func GetErrorBase(e error) Error {
	return e.(Error)
}

func IsErrorOfType(e error, signatureError string) bool {
	if !IsErrorBase(e) {
		return false
	}

	eN := GetErrorBase(e)
	return eN.GetSignature() == signatureError
}

func New(code int, message string) Error {
	return newBaseError(code, message)
}

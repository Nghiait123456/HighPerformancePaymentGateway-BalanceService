package auth_internal_service

import "time"

type (
	AuthInternalService struct {
		secretNearestOld string
		secretCurrent    string
		timeUpdateSecret uint // ms
	}

	AuthResponse struct {
		Status  string
		Code    string
		Message string
	}

	AuthInternalServiceInterface interface {
		GetSecret() string                            // secret use auth with other service
		Auth(secretCheck string) (bool, AuthResponse) // auth to incoming request
		Init(timeUpdateSecret uint) error
		autoUpdateSecret()
		UpdateSecret() error
	}
)

const (
	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"

	CODE_SUCCESS    = "SUCCESS"
	MESSAGE_SUCCESS = "auth success"

	CODE_UPDATE_NEW_SECRET    = "UPDATE_NEW_SECRET"
	MESSAGE_UPDATE_NEW_SECRET = "auth success, but have new secret source,  please update new secret"

	CODE_ERROR    = "ERROR"
	MESSAGE_ERROR = "secret not correct, please update new secret and try again"
)

func (a *AuthInternalService) Auth(secretCheck string) (bool, AuthResponse) {
	if a.secretCurrent == secretCheck {
		return true, AuthResponse{
			Status:  STATUS_SUCCESS,
			Code:    CODE_SUCCESS,
			Message: MESSAGE_SUCCESS,
		}
	}

	if a.secretNearestOld == secretCheck {
		return true, AuthResponse{
			Status:  STATUS_SUCCESS,
			Code:    CODE_UPDATE_NEW_SECRET,
			Message: MESSAGE_UPDATE_NEW_SECRET,
		}
	}

	return false, AuthResponse{
		Status:  STATUS_ERROR,
		Code:    CODE_ERROR,
		Message: MESSAGE_ERROR,
	}
}

func (a *AuthInternalService) GetSecret() string {
	return a.secretCurrent
}

func (a *AuthInternalService) UpdateSecret() error {
	// todo call api Update Secret
	// todo if secret get != currenct secret, update  secretNearestOld and secretCurrent
	return nil
}

func (a *AuthInternalService) autoUpdateSecret() {
	for {
		err := a.UpdateSecret()
		if err != nil {
			//todo log
		}

		// Sleep() is the best for optimal cpu resource
		time.Sleep(time.Duration(a.timeUpdateSecret) * time.Millisecond)
	}
}

func (a *AuthInternalService) Init(timeUpdateSecret uint) error {
	a.timeUpdateSecret = timeUpdateSecret
	go func() {
		a.autoUpdateSecret()
	}()

	return nil
}

func NewAuthInternalService() AuthInternalServiceInterface {
	return &AuthInternalService{}
}

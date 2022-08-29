package auth_internal_service

const (
	HTTP_CODE_SUCCESS = 200
	HTTP_CODE_ERROR   = 401

	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"

	CODE_SUCCESS    = 0
	MESSAGE_SUCCESS = "auth success"

	DEFAULT_ERROR_CODE    = 501
	DEFAULT_MESSAGE_ERROR = "Unknow error code, please update new secret and try again"

	CODE_ERROR_VERSION_OF_SECRET_WRONG_FORMAT    = 1
	MESSAGE_ERROR_VERSION_OF_SECRET_WRONG_FORMAT = "wrong format version of secret code, version must is uint and length is 20"

	CODE_UPDATE_NEW_SECRET    = 2
	MESSAGE_UPDATE_NEW_SECRET = "auth success, but have new secret source,  please update new secret"

	CODE_ERROR_VERSION_OF_SECRET_OLD    = 3
	MESSAGE_ERROR_VERSION_OF_SECRET_OLD = "Version of secret is old, please update new secret and try again"

	CODE_ERROR_WRONG_SECRET    = 4
	MESSAGE_ERROR_WRONG_SECRET = "secret not correct, please update new secret and try again"
)

func responseAuthSuccess() (bool, AuthResponse) {
	return true, AuthResponse{
		HttpCode: HTTP_CODE_SUCCESS,
		Status:   STATUS_SUCCESS,
		Code:     CODE_SUCCESS,
		Message:  MESSAGE_SUCCESS,
	}
}

func responseAuthSuccessButUpdateNewSecret() (bool, AuthResponse) {
	return true, AuthResponse{
		HttpCode: HTTP_CODE_SUCCESS,
		Status:   STATUS_SUCCESS,
		Code:     CODE_UPDATE_NEW_SECRET,
		Message:  MESSAGE_UPDATE_NEW_SECRET,
	}
}

func responseAuthErrorWhyVersionSecretKeyWrongFormat() (bool, AuthResponse) {
	return false, AuthResponse{
		HttpCode: HTTP_CODE_ERROR,
		Status:   STATUS_ERROR,
		Code:     CODE_ERROR_VERSION_OF_SECRET_WRONG_FORMAT,
		Message:  MESSAGE_ERROR_VERSION_OF_SECRET_WRONG_FORMAT,
	}
}

func responseAuthErrorWhyVersionSecretOld() (bool, AuthResponse) {
	return false, AuthResponse{
		HttpCode: HTTP_CODE_ERROR,
		Status:   STATUS_ERROR,
		Code:     CODE_ERROR_VERSION_OF_SECRET_OLD,
		Message:  MESSAGE_ERROR_VERSION_OF_SECRET_OLD,
	}
}

func defaultErrorWhyNotClearReason() (bool, AuthResponse) {
	return false, AuthResponse{
		HttpCode: HTTP_CODE_ERROR,
		Status:   STATUS_ERROR,
		Code:     DEFAULT_ERROR_CODE,
		Message:  DEFAULT_MESSAGE_ERROR,
	}
}

func responseAuthErrorWhyWrongSecret() (bool, AuthResponse) {
	return false, AuthResponse{
		HttpCode: HTTP_CODE_ERROR,
		Status:   STATUS_ERROR,
		Code:     CODE_ERROR_WRONG_SECRET,
		Message:  MESSAGE_ERROR_WRONG_SECRET,
	}
}

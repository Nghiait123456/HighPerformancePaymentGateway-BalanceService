package auth_internal_service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type (
	AutoChangeSecret struct {
	}

	AutoChangeSecretInterface interface {
		ReplaceNewSecret()
		CreateNewSecret(currentSecret string) (string, bool)
		ValidateSecret(secret string) (bool, error)
	}
)

func (a *AutoChangeSecret) ValidateSecret(secret string) (bool, error) {
	if len(secret) != LENGTH_SECRET_STRING {
		log.WithFields(log.Fields{
			"lengSecret":   len(secret),
			"errorMessage": errV.Error(),
		}).Errorf("Wrong leng secret, must is %d", LENGTH_SECRET_STRING)
		panic(fmt.Sprintf("Wrong leng secret, must is %d", LENGTH_SECRET_STRING))
	}

	return true, nil
}

func (a *AutoChangeSecret) CreateNewSecret(currentSecret string) (string, bool) {

}

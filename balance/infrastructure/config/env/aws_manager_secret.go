package env

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	AwsMangerSecret struct {
		SecretName   string
		Region       string
		VersionState string
	}

	AwsMangerSecretInterface interface {
		Init(secretName string, region string, versionState string)
		GetSecret() (string, error)             // json string
		UpdateSecretString(update string) error // update secret type json string : { key : value, ...}
		IsErrorCustomOfAws(e error) bool
		ConvertToErrorAws(e error) awserr.Error
	}
)

const (
	DEFAULT_AWS_SERVSION_STATTE = "AWSCURRENT"
)

func (a *AwsMangerSecret) Init(secretName string, region string, versionState string) {
	if secretName == "" {
		panic("secretName must not empty")
		os.Exit(0)
	}

	if region == "" {
		panic("region must not empty")
		os.Exit(0)
	}

	if versionState == "" {
		versionState = DEFAULT_AWS_SERVSION_STATTE
	}

	a.SecretName = secretName
	a.Region = region
	a.VersionState = versionState
}

func (a AwsMangerSecret) GetSecret() (string, error) {
	sess, errNS := session.NewSession()
	if errNS != nil {
		return "", errNS
	}

	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(a.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(a.SecretName),
		VersionStage: aws.String(a.VersionState),
	}

	result, errG := svc.GetSecretValue(input)
	if errG != nil {
		if aerr, ok := errG.(awserr.Error); ok {
			log.WithFields(log.Fields{
				"errCode":    aerr.Code(),
				"errMessage": aerr.Message(),
			}).Error("Update secret error")

			return "", aerr
		} else {
			log.WithFields(log.Fields{
				"errMessage": errG.Error(),
			}).Error("Update secret error")

			return "", errG
		}
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		secretString := *result.SecretString
		return secretString, nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return "", err
		}

		decodedBinarySecret := string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret, nil
	}
}

// update is json string
func (a AwsMangerSecret) UpdateSecretString(update string) error {
	sess, errNS := session.NewSession()
	if errNS != nil {
		return errNS
	}

	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(a.Region))
	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(a.SecretName),
		SecretString: aws.String(update),
	}

	_, errUpdate := svc.UpdateSecret(input)
	if errUpdate != nil {
		if aerr, ok := errUpdate.(awserr.Error); ok {
			log.WithFields(log.Fields{
				"errCode":    aerr.Code(),
				"errMessage": aerr.Message(),
			}).Error("Update secret error")
		} else {
			log.WithFields(log.Fields{
				"errMessage": errUpdate.Error(),
			}).Error("Update secret error")
		}

		return errUpdate
	}

	return nil
}

func (a AwsMangerSecret) IsErrorCustomOfAws(e error) bool {
	if _, ok := e.(awserr.Error); ok {
		return true
	}

	return false
}

func (a AwsMangerSecret) ConvertToErrorAws(e error) awserr.Error {
	return e.(awserr.Error)
}

func NewAwsManagerSecret() AwsMangerSecretInterface {
	return &AwsMangerSecret{}
}

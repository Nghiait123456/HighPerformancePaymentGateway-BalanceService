package env

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type (
	AwsMangerSecret struct {
		SecretName   string
		Region       string
		VersionState string
	}

	AwsMangerSecretInterface interface {
		Init(secretName string, region string, versionState string)
		GetSecret() (string, error) // json string
	}
)

func (a *AwsMangerSecret) Init(secretName string, region string, versionState string) {
	a.SecretName = secretName
	a.Region = region
	a.VersionState = versionState
}

func (a AwsMangerSecret) GetSecret() (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(a.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(a.SecretName),
		VersionStage: aws.String(a.VersionState),
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return "", aerr
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				//todo log error
				return "", aerr

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				//todo log error
				return "", aerr

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				//todo log error
				return "", aerr

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				//todo log error
				return "", aerr

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				//todo log error
				return "", aerr
			}
		} else {
			fmt.Println(err.Error())
			//todo log error
			return "", aerr
		}
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	if result.SecretString != nil {
		secretString := *result.SecretString
		fmt.Println(fmt.Sprintf("secret string %s", secretString))
		return secretString, nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			//todo log error
			return "", err
		}

		decodedBinarySecret := string(decodedBinarySecretBytes[:len])
		fmt.Println(fmt.Sprintf("decodedBinarySecretg %s", decodedBinarySecret))
		return decodedBinarySecret, nil
	}
}

func NewAwsManagerSecret() AwsMangerSecretInterface {
	return &AwsMangerSecret{}
}

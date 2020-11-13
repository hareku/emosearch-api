package secrets

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var sm *secretsmanager.SecretsManager

func init() {
	sm = secretsmanager.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))
}

// Get gets the specified key value from Amazon Secrets Manager.
func Get(secretID string) (*string, error) {
	result, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretID),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return nil, fmt.Errorf("aws secrets manager error: %w", err)
	}

	return result.SecretString, nil
}

// GetBase64 gets the specified base64 encoded value from Amazon Secrets Manager and decodes it.
func GetBase64(secretID string) ([]byte, error) {
	raw, err := Get(secretID)
	if err != nil {
		return nil, fmt.Errorf("failed to get base64 encoded value from aws secrets manager: %w", err)
	}

	val, err := base64.StdEncoding.DecodeString(*raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 value: %w", err)
	}

	return val, nil
}

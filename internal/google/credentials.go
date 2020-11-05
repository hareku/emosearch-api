package google

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// GetGoogleServiceAccountKey returns Google service account key.
func GetGoogleServiceAccountKey() ([]byte, error) {
	envKey := getKeyFromEnv()
	if envKey != "" {
		decodedString, err := base64.StdEncoding.DecodeString(envKey)

		if err != nil {
			return nil, fmt.Errorf("base64 decoding failed: %w", err)
		}

		return decodedString, nil
	}

	smsKey, err := getKeyFromSMS()
	if err != nil {
		return nil, fmt.Errorf("cloudn't get key from sms: %w", err)
	}

	decodedString, err := base64.StdEncoding.DecodeString(smsKey)

	if err != nil {
		return nil, fmt.Errorf("base64 decoding failed: %w", err)
	}

	return decodedString, nil
}

func getKeyFromEnv() string {
	return os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY")
}

func getKeyFromSMS() (string, error) {
	secretName := "GoogleServiceAccountKey"
	region := "ap-northeast-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("aws sms error: %w", err)
	}

	return *result.SecretString, nil
}

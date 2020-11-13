package google

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/hareku/emosearch-api/internal/secrets"
)

// GetGoogleServiceAccountKey returns Google service account key.
func GetGoogleServiceAccountKey() ([]byte, error) {
	// First, we try to get the key from env.
	envVal := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY")
	if envVal != "" {
		decoded, err := base64.StdEncoding.DecodeString(envVal)
		if err != nil {
			return nil, fmt.Errorf("failed to decode GOOGLE_SERVICE_ACCOUNT_KEY env as base64: %w", err)
		}
		return decoded, nil
	}

	// Next, we try to get the key from Amazon Secrets Manager.
	smArn := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_SECRETS_MANAGER_ARN")
	if smArn != "" {
		smVal, err := secrets.Get(smArn)
		if err != nil {
			return nil, fmt.Errorf("failed to get %w", err)
		}
		if smVal != nil {
			decoded, err := base64.StdEncoding.DecodeString(*smVal)
			if err != nil {
				return nil, fmt.Errorf("failed to decode key as base64 which is stored aws %s: %w", smArn, err)
			}
			return decoded, nil
		}
	}

	return nil, errors.New("google service account key was not found")
}

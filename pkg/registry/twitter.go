package registry

import (
	"errors"
	"fmt"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/hareku/emosearch-api/internal/secrets"
	domain_twitter "github.com/hareku/emosearch-api/pkg/domain/twitter"
	infra_twitter "github.com/hareku/emosearch-api/pkg/infrastructure/twitter"
)

var oauth1Config *oauth1.Config

func getOauth1Config() *oauth1.Config {
	if oauth1Config == nil {
		key, err := getTwitterConsumerKey()
		if err != nil {
			panic(fmt.Errorf("twitter oauth1 client initialization error: %w", err))
		}
		secret, err := getTwitterConsumerSecret()
		if err != nil {
			panic(fmt.Errorf("twitter oauth1 client initialization error: %w", err))
		}

		oauth1Config = oauth1.NewConfig(*key, *secret)
	}
	return oauth1Config
}

func getTwitterConsumerKey() (*string, error) {
	// First, we try to get the key from env.
	envVal := os.Getenv("TWITTER_CONSUMER_KEY")
	if envVal != "" {
		return &envVal, nil
	}

	// Next, we try to get the key from Amazon Secrets Manager.
	smArn := os.Getenv("TWITTER_CONSUMER_KEY_SECRETS_MANAGER_ARN")
	if smArn != "" {
		smVal, err := secrets.Get(smArn)
		if err != nil {
			return nil, fmt.Errorf("failed to get twitter consumer key (%s) from secrets manager: %w", smArn, err)
		}
		return smVal, nil
	}

	return nil, errors.New("twitter consumer key was not found")
}

func getTwitterConsumerSecret() (*string, error) {
	// First, we try to get the key from env.
	envVal := os.Getenv("TWITTER_CONSUMER_SECRET")
	if envVal != "" {
		return &envVal, nil
	}

	// Next, we try to get the key from Amazon Secrets Manager.
	smArn := os.Getenv("TWITTER_CONSUMER_SECRET_SECRETS_MANAGER_ARN")
	if smArn != "" {
		smVal, err := secrets.Get(smArn)
		if err != nil {
			return nil, fmt.Errorf("failed to get twitter consumer secret (%s) from secrets manager: %w", smArn, err)
		}
		return smVal, nil
	}

	return nil, errors.New("twitter consumer secret was not found")
}

// NewTwitterClient create Client of domain Twitter.
func (r *registry) NewTwitterClient() domain_twitter.Client {
	return infra_twitter.NewTwitterOauth1Client(getOauth1Config())
}

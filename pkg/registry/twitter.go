package registry

import (
	"os"

	"github.com/dghubble/oauth1"
	domain_twitter "github.com/hareku/emosearch-api/pkg/domain/twitter"
	infra_twitter "github.com/hareku/emosearch-api/pkg/infrastructure/twitter"
)

var oauth1Config *oauth1.Config

func getOauth1Config() *oauth1.Config {
	if oauth1Config == nil {
		oauth1Config = oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	}
	return oauth1Config
}

// NewTwitterClient create Client of domain Twitter.
func (r *registry) NewTwitterClient() domain_twitter.Client {
	return infra_twitter.NewTwitterOauth1Client(getOauth1Config())
}

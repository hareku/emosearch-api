package repository

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/model"
)

// TweetRepository provides CRUD methods for Tweet domain.
type TweetRepository interface {
	Store(ctx context.Context, tweet *model.Tweet) error
}

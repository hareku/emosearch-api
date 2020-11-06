package dynamodb

import (
	"context"
	"fmt"

	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

type dynamoDbUserRepository struct {
	dynamoDB dynamo.Table
}

// NewDynamoDatabaseUserRepository creates UserRepository which implemented by DynamoDB.
func NewDynamoDatabaseUserRepository(dynamoDB dynamo.Table) repository.UserRepository {
	return &dynamoDbUserRepository{dynamoDB}
}

type dynamoDBUser struct {
	PK string
	SK string

	UserID                   model.UserID `dynamo:"UserID"`
	TwitterAccessToken       string       `dynamo:"TwitterAccessToken"`
	TwitterAccessTokenSecret string       `dynamo:"TwitterAccessTokenSecret"`
}

func (r *dynamoDbUserRepository) Create(ctx context.Context, user *model.User) error {
	_user := dynamoDBUser{
		PK: fmt.Sprintf("USER#%s", user.UserID),
		SK: fmt.Sprintf("PROFILE#%s", user.UserID),

		UserID:                   user.UserID,
		TwitterAccessToken:       user.TwitterAccessToken,
		TwitterAccessTokenSecret: user.TwitterAccessTokenSecret,
	}

	err := r.dynamoDB.Put(&_user).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("DynamoDB error: %w", err)
	}

	return nil
}

func (r *dynamoDbUserRepository) FindByID(ctx context.Context, userID model.UserID) (*model.User, error) {
	dbUser := &dynamoDBUser{}

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Range("SK", dynamo.Equal, fmt.Sprintf("PROFILE#%s", userID)).
		OneWithContext(ctx, dbUser)

	if err != nil {
		return nil, fmt.Errorf("DynamoDB error: %w", err)
	}

	user := &model.User{
		UserID:                   dbUser.UserID,
		TwitterAccessToken:       dbUser.TwitterAccessToken,
		TwitterAccessTokenSecret: dbUser.TwitterAccessTokenSecret,
	}

	return user, nil
}

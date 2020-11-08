package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/internal/uuid"
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

type dynamoDBSearchRepository struct {
	dynamoDB dynamo.Table
}

// NewDynamoDBSearchRepository creates SearchRepository which is implemented by DynamoDB.
func NewDynamoDBSearchRepository(dynamoDB dynamo.Table) repository.SearchRepository {
	return &dynamoDBSearchRepository{dynamoDB}
}

type dynamoDBSearch struct {
	PK string
	SK string

	SearchID  model.SearchID             `dynamo:"SearchID"`
	UserID    model.UserID               `dynamo:"UserID"`
	Title     string                     `dynamo:"Title"`
	Query     string                     `dynamo:"Query"`
	CreatedAt dynamodbattribute.UnixTime `dynamo:"CreatedAt"`
	UpdatedAt dynamodbattribute.UnixTime `dynamo:"UpdatedAt"`
}

func (d *dynamoDBSearch) NewSearchModel() *model.Search {
	return &model.Search{
		SearchID:  d.SearchID,
		UserID:    d.UserID,
		Title:     d.Title,
		Query:     d.Query,
		CreatedAt: time.Time(d.CreatedAt),
		UpdatedAt: time.Time(d.UpdatedAt),
	}
}

func (r *dynamoDBSearchRepository) ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error) {
	var result []dynamoDBSearch
	var searches []*model.Search

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Range("SK", dynamo.BeginsWith, "SEARCH#").
		AllWithContext(ctx, &result)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return searches, nil
			default:
				return nil, fmt.Errorf("DynamoDB error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("DynamoDB error: %w", err)
		}
	}

	for i := 0; i < len(result); i++ {
		searches = append(searches, result[i].NewSearchModel())
	}

	if len(searches) == 0 {
		searches = []*model.Search{}
	}

	return searches, nil
}

func (r *dynamoDBSearchRepository) Create(ctx context.Context, search *model.Search) error {
	searchID, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("uuid error: %w", err)
	}

	search.SearchID = model.SearchID(searchID)

	createdAt := time.Now()
	search.CreatedAt = createdAt
	search.UpdatedAt = createdAt

	dsearch := dynamoDBSearch{
		PK: fmt.Sprintf("USER#%s", search.UserID),
		SK: fmt.Sprintf("SEARCH#%s", search.SearchID),

		UserID:    search.UserID,
		SearchID:  search.SearchID,
		Title:     search.Title,
		Query:     search.Query,
		CreatedAt: dynamodbattribute.UnixTime(search.CreatedAt),
		UpdatedAt: dynamodbattribute.UnixTime(search.UpdatedAt),
	}

	err = r.dynamoDB.Put(&dsearch).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("DynamoDB error: %w", err)
	}

	return nil
}

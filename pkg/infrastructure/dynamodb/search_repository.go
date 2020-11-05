package dynamodb

import (
	"context"
	"fmt"

	"github.com/guregu/dynamo"
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

	UserID    string `dynamo:"UserID"`
	SearchID  string `dynamo:"SearchID"`
	Title     string `dynamo:"Title"`
	Query     string `dynamo:"Query"`
	CreatedAt string `dynamo:"CreatedAt"`
	UpdatedAt string `dynamo:"UpdatedAt"`
}

func (d *dynamoDBSearch) NewSearchModel() *model.Search {
	return &model.Search{
		UserID:    d.UserID,
		SearchID:  d.SearchID,
		Title:     d.Title,
		Query:     d.Query,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *dynamoDBSearchRepository) ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error) {
	var result []dynamoDBSearch

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Filter("BEGINS_WITH(SK, 'SEARCH#')").
		AllWithContext(ctx, &result)

	if err != nil {
		return nil, fmt.Errorf("DynamoDB error: %w", err)
	}

	var searches []*model.Search

	for i := 0; i < len(result); i++ {
		searches = append(searches, result[i].NewSearchModel())
	}

	return searches, nil
}

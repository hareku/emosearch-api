package dynamodb

import (
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

func (r *dynamoDBSearchRepository) ListByUserID(userID model.UserID) ([]model.Search, error) {
	var result []model.Search

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Filter("BEGINS_WITH(SK, 'SEARCH#')").
		All(&result)

	if err != nil {
		return nil, fmt.Errorf("DynamoDB error: %w", err)
	}

	return result, nil
}

package dynamodb

import (
	"context"
	"errors"
	"fmt"

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

// This value is used for "SearchIndex" GSI of DynamoDB.
// May be split into random numbers in the future.
const searchIndexPK = 1

type dynamoDBSearch struct {
	PK            string
	SK            string
	SearchIndexPK int64
	*model.Search
}

func (d *dynamoDBSearch) NewSearchModel() *model.Search {
	return d.Search
}

func (r *dynamoDBSearchRepository) List(ctx context.Context, input repository.SearchRepositoryListInput) ([]*model.Search, error) {
	var items []dynamoDBSearch

	q := r.dynamoDB.Get("SearchIndexPK", searchIndexPK).
		Index("SearchIndex").
		Order(false).
		Limit(input.Limit)

	if input.UntilNextSearchUpdateAt != nil {
		q.Range("NextSearchUpdateAt", dynamo.LessOrEqual, *input.UntilNextSearchUpdateAt)
	}

	err := q.All(&items)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamo error: %w", err)
	}

	var res []*model.Search
	for _, item := range items {
		res = append(res, item.NewSearchModel())
	}
	return res, nil
}

func (r *dynamoDBSearchRepository) ListByUserID(ctx context.Context, userID model.UserID) ([]*model.Search, error) {
	var dynamoResult []dynamoDBSearch
	searches := []*model.Search{}

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Range("SK", dynamo.BeginsWith, "SEARCH#").
		AllWithContext(ctx, &dynamoResult)

	if errors.Is(err, dynamo.ErrNotFound) {
		return searches, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamo error: %w", err)
	}

	for i := 0; i < len(dynamoResult); i++ {
		searches = append(searches, dynamoResult[i].NewSearchModel())
	}

	return searches, nil
}

func (r *dynamoDBSearchRepository) Find(ctx context.Context, userID model.UserID, searchID model.SearchID) (*model.Search, error) {
	var dynamoSearch dynamoDBSearch

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Range("SK", dynamo.Equal, fmt.Sprintf("SEARCH#%s", searchID)).
		OneWithContext(ctx, &dynamoSearch)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamo error: %w", err)
	}

	return dynamoSearch.NewSearchModel(), nil
}

func (r *dynamoDBSearchRepository) Create(ctx context.Context, search *model.Search) error {
	searchID, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("uuid error: %w", err)
	}
	search.SearchID = model.SearchID(searchID)

	dynamoSearch := dynamoDBSearch{
		PK:            fmt.Sprintf("USER#%s", search.UserID),
		SK:            fmt.Sprintf("SEARCH#%s", search.SearchID),
		SearchIndexPK: searchIndexPK,
		Search:        search,
	}

	err = r.dynamoDB.Put(&dynamoSearch).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}

func (r *dynamoDBSearchRepository) Update(ctx context.Context, search *model.Search) error {
	err := r.dynamoDB.Update("PK", fmt.Sprintf("USER#%s", search.UserID)).
		Range("SK", fmt.Sprintf("SEARCH#%s", search.SearchID)).
		Set("Query", search.Query).
		Set("NextSearchUpdateAt", search.NextSearchUpdateAt).
		RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}

func (r *dynamoDBSearchRepository) Delete(ctx context.Context, search *model.Search) error {
	err := r.dynamoDB.Delete("PK", fmt.Sprintf("USER#%s", search.UserID)).
		Range("SK", fmt.Sprintf("SEARCH#%s", search.SearchID)).
		RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("dynamo error: %w", err)
	}

	return nil
}

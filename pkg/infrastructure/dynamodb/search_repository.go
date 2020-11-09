package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"time"

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

	SearchIndexID model.SearchID // for global secondery index "SearchIndex"
	SearchID      model.SearchID
	UserID        model.UserID
	Title         string
	Query         string
	LastUpdatedAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	var dSearch dynamoDBSearch

	err := r.dynamoDB.
		Get("PK", fmt.Sprintf("USER#%s", userID)).
		Range("SK", dynamo.Equal, fmt.Sprintf("SEARCH#%s", searchID)).
		OneWithContext(ctx, &dSearch)

	if errors.Is(err, dynamo.ErrNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("dynamo error: %w", err)
	}

	return dSearch.NewSearchModel(), nil
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

		UserID:        search.UserID,
		SearchIndexID: search.SearchID,
		SearchID:      search.SearchID,
		LastUpdatedAt: createdAt.AddDate(-1, 0, 0),
		Title:         search.Title,
		Query:         search.Query,
		CreatedAt:     search.CreatedAt,
		UpdatedAt:     search.UpdatedAt,
	}

	err = r.dynamoDB.Put(&dsearch).RunWithContext(ctx)

	if err != nil {
		return fmt.Errorf("DynamoDB error: %w", err)
	}

	return nil
}

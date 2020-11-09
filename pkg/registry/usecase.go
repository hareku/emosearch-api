package registry

import (
	"github.com/hareku/emosearch-api/pkg/usecase"
)

func (r *registry) NewUserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(r.NewAuthenticator(), r.NewUserRepository())
}

func (r *registry) NewSearchUsecase() usecase.SearchUsecase {
	return usecase.NewSearchUsecase(r.NewAuthenticator(), r.NewValidator(), r.NewSearchRepository())
}

func (r *registry) NewBatchUsecase() usecase.BatchUsecase {
	return usecase.NewBatchUsecase(&usecase.NewBatchUsecaseInput{
		Authenticator:     r.NewAuthenticator(),
		UserUsecase:       r.NewUserUsecase(),
		SearchUsecase:     r.NewSearchUsecase(),
		TweetRepository:   r.NewTweetRepository(),
		TwitterClient:     r.NewTwitterClient(),
		SentimentDetector: r.NewSentimentDetector(),
	})
}

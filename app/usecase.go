package app

import (
	"context"

	"github.com/tomocy/mint/domain"
)

func NewTweetUsecase(repo domain.TweetRepository) *TweetUsecase {
	return &TweetUsecase{
		repo: repo,
	}
}

type TweetUsecase struct {
	repo domain.TweetRepository
}

func (u *TweetUsecase) PoleHomeTweets(ctx context.Context) (<-chan []*domain.Tweet, <-chan error) {
	return u.repo.PoleHomeTweets(ctx)
}

func (u *TweetUsecase) FetchHomeTweets() ([]*domain.Tweet, error) {
	return u.repo.FetchHomeTweets()
}

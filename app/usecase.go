package app

import (
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

func (u *TweetUsecase) Fetch() ([]*domain.Tweet, error) {
	return u.repo.FetchHomeTweets()
}

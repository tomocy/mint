package domain

import "context"

type TweetRepository interface {
	PoleHomeTweets(ctx context.Context) (<-chan []*Tweet, <-chan error)
	FetchHomeTweets() ([]*Tweet, error)
}

package domain

type TweetRepository interface {
	FetchHomeTweets() ([]*Tweet, error)
}

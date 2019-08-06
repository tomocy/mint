package client

import (
	"fmt"

	"github.com/tomocy/mint/app"
	"github.com/tomocy/mint/domain"
	"github.com/tomocy/mint/infra"
)

type CLI struct{}

func (c *CLI) FetchHomeTweets() error {
	r := infra.NewTwitter()
	u := app.NewTweetUsecase(r)
	tweets, err := u.FetchHomeTweets()
	if err != nil {
		return err
	}

	c.showTweets(tweets)

	return nil
}

func (c *CLI) showTweets(tweets []*domain.Tweet) {
	stringer := asciiTweets(tweets)
	fmt.Print(stringer)
}

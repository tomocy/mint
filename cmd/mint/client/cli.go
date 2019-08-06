package client

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomocy/mint/app"
	"github.com/tomocy/mint/domain"
	"github.com/tomocy/mint/infra"
)

type CLI struct{}

func (c *CLI) PoleHomeTweets() error {
	r := infra.NewTwitter()
	u := app.NewTweetUsecase(r)
	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)
	tsCh, errCh := u.PoleHomeTweets(ctx)
	for {
		select {
		case ts := <-tsCh:
			ordered := orderOlderTweets(ts)
			c.showTweets(ordered)
			fmt.Printf("updated at %s\n", time.Now().Format("2006/01/02 15:04"))
		case err := <-errCh:
			cancelFn()
			return err
		case sig := <-sigCh:
			cancelFn()
			fmt.Println(sig)
			return nil
		}
	}
}

func (c *CLI) FetchHomeTweets() error {
	r := infra.NewTwitter()
	u := app.NewTweetUsecase(r)
	tweets, err := u.FetchHomeTweets()
	if err != nil {
		return err
	}

	ordered := orderOlderTweets(tweets)
	c.showTweets(ordered)

	return nil
}

func (c *CLI) showTweets(tweets []*domain.Tweet) {
	stringer := asciiTweets(tweets)
	fmt.Print(stringer)
}

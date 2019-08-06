package client

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tomocy/mint/domain"
)

type asciiTweets []*domain.Tweet

func (ts asciiTweets) String() string {
	var b strings.Builder
	for i, t := range ts {
		if i == 0 {
			b.WriteString("----------\n")
		}
		b.WriteString(fmt.Sprintf("%s@%s %s\n%s\n", t.User.Name, t.User.ScreenName, t.CreatedAt.Format("2006/01/02 15:04"), t.Text))
		b.WriteString("----------\n")
	}

	return b.String()
}

func orderOlderTweets(ts []*domain.Tweet) []*domain.Tweet {
	ordered := make([]*domain.Tweet, len(ts))
	copy(ordered, ts)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].CreatedAt.Before(ordered[j].CreatedAt)
	})

	return ordered
}

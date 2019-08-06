package client

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tomocy/mint/domain"
)

type asciiTweets []*domain.Tweet

func (ts asciiTweets) String() string {
	var b strings.Builder
	for i, t := range ts {
		if i == 0 {
			b.WriteString("----------\n")
		}
		b.WriteString(fmt.Sprintf("%s@%s %s\n%s\n", t.User.Name, t.User.ScreenName, legibleTimeDiff(t.CreatedAt, "2006/01/02 15:04"), t.Text))
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

func legibleTimeDiff(t time.Time, format string) string {
	now := time.Now()
	if now.Add(-7 * 24 * time.Hour).After(t) {
		return t.Format(format)
	}

	diff := now.Sub(t)
	if 24 <= diff.Hours() {
		return fmt.Sprintf("%dd", int(diff.Hours()/24))
	}
	if 1 <= diff.Hours() {
		return fmt.Sprintf("%dh", int(diff.Hours()))
	}
	if 1 <= diff.Minutes() {
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	}

	return fmt.Sprintf("%ds", int(diff.Seconds()))
}

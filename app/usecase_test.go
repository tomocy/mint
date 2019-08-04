package app

import (
	"fmt"
	"testing"

	"github.com/tomocy/mint/domain"
)

func TestFetch(t *testing.T) {
	testers := map[string]func(t *testing.T){
		"tweet usecase": testTweetUsecaseFetch,
	}

	for name, tester := range testers {
		t.Run(name, tester)
	}
}

func testTweetUsecaseFetch(t *testing.T) {
	r := new(mock)
	u := NewTweetUsecase(r)
	expecteds := []*domain.Tweet{}
	actuals, err := u.Fetch()
	if err != nil {
		t.Errorf("unexpected error returned by Fetch: got %s, expect nil\n", err)
	}
	if len(actuals) != len(expecteds) {
		t.Fatalf("unexpected length of tweets returned by Fetch: got %d, expected %d\n", len(actuals), len(expecteds))
	}
	for i, expected := range expecteds {
		if err := assertTweet(actuals[i], expected); err != nil {
			t.Errorf("unexpected tweet returned by Fetch: %s\n", err)
		}
	}
}

func assertTweet(actual, expected *domain.Tweet) error {
	if actual.ID != expected.ID {
		return reportUnexpected("id of tweet", actual.ID, expected.ID)
	}
	if actual.Text != expected.Text {
		return reportUnexpected("text of tweet", actual.Text, expected.Text)
	}
	if err := assertUser(actual.User, expected.User); err != nil {
		return fmt.Errorf("unexpected user or tweet: %s", err)
	}
	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		return reportUnexpected("created at of tweet", actual.CreatedAt, expected.CreatedAt)
	}

	return nil
}

func assertUser(actual, expected *domain.User) error {
	if actual.ID != expected.ID {
		return reportUnexpected("id of user", actual.ID, expected.ID)
	}
	if actual.Name != expected.Name {
		return reportUnexpected("name of user", actual.Name, expected.Name)
	}
	if actual.ScreenName != expected.ScreenName {
		return reportUnexpected("screen name of user", actual.ScreenName, expected.ScreenName)
	}

	return nil
}

func reportUnexpected(name string, actual, expected interface{}) error {
	return fmt.Errorf("unexpected %s: got %v, expect %v", name, actual, expected)
}

type mock struct{}

func (m *mock) FetchHomeTweets() ([]*domain.Tweet, error) {
	return []*domain.Tweet{}, nil
}

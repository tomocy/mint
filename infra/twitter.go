package infra

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/tomocy/mint/domain"
	"github.com/tomocy/mint/infra/twitter"
)

func NewTwitter() *Twitter {
	createWorkspace()
	return &Twitter{
		oauthClient: &oauth.Client{
			TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
			ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
			TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
			Credentials: oauth.Credentials{
				Token:  "HW2BeqW5CiP2jJW9KzpWu7GIo",
				Secret: "zkhkGA72sWK2Arsc1gt4XcjDeD9Ispfx1gS4a0sC6wSlT8Vl2u",
			},
		},
	}
}

type Twitter struct {
	oauthClient *oauth.Client
}

func (t *Twitter) PoleHomeTweets(ctx context.Context) (<-chan []*domain.Tweet, <-chan error) {
	tsCh, errCh := make(chan []*domain.Tweet), make(chan error)
	go func() {
		defer func() {
			close(tsCh)
			close(errCh)
		}()

		sendTweets := func(lastID string, tsCh chan<- []*domain.Tweet, errCh chan<- error) string {
			var params url.Values
			if lastID != "" {
				params = url.Values{
					"since_id": []string{lastID},
				}
			}
			ts, err := t.fetchTweets("https://api.twitter.com/1.1/statuses/home_timeline.json", params)
			if err != nil {
				errCh <- err
				return ""
			}
			if len(ts) <= 0 {
				return ""
			}

			tsCh <- ts
			return ts[0].ID
		}

		lastID := sendTweets("", tsCh, errCh)
		for {
			select {
			case <-ctx.Done():
				break
			case <-time.After(4 * time.Minute):
				lastID = sendTweets(lastID, tsCh, errCh)
			}
		}
	}()

	return tsCh, errCh
}

func (t *Twitter) poleTweets(ctx context.Context, rawURL string) (<-chan []*domain.Tweet, <-chan error) {
	tsCh, errCh := make(chan []*domain.Tweet), make(chan error)
	go func() {
		defer func() {
			close(tsCh)
			close(errCh)
		}()

		sendTweets := func(lastID string, tsCh chan<- []*domain.Tweet, errCh chan<- error) string {
			var params url.Values
			if lastID != "" {
				params = url.Values{
					"since_id": []string{lastID},
				}
			}
			ts, err := t.fetchTweets(rawURL, params)
			if err != nil {
				errCh <- err
				return ""
			}
			if len(ts) <= 0 {
				return ""
			}

			tsCh <- ts
			return ts[0].ID
		}

		lastID := sendTweets("", tsCh, errCh)
		for {
			select {
			case <-ctx.Done():
				break
			case <-time.After(4 * time.Minute):
				lastID = sendTweets(lastID, tsCh, errCh)
			}
		}
	}()

	return tsCh, errCh
}

func (t *Twitter) FetchHomeTweets() ([]*domain.Tweet, error) {
	ts, err := t.fetchTweets("https://api.twitter.com/1.1/statuses/home_timeline.json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch home tweets: %s", err)
	}

	return ts, nil
}

func (t *Twitter) fetchTweets(url string, params url.Values) ([]*domain.Tweet, error) {
	cred, err := t.treiveAuthorization()
	if err != nil {
		return nil, err
	}

	var tweets twitter.Tweets
	if err := t.trieve(&oauthRequest{
		cred:   cred,
		method: http.MethodGet,
		url:    url,
		params: params,
	}, &tweets); err != nil {
		t.deleteConfig()
		return nil, err
	}

	if err := t.saveConfig(&config{
		AccessCredentials: cred,
	}); err != nil {
		return nil, err
	}

	return tweets.Adapt(), nil
}

func (t *Twitter) deleteConfig() {
	os.Remove(configFilename())
}

func (t *Twitter) saveConfig(config *config) error {
	destName := configFilename()
	dest, err := os.OpenFile(destName, os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer dest.Close()

	return json.NewEncoder(dest).Encode(config)
}

func (t *Twitter) treiveAuthorization() (*oauth.Credentials, error) {
	config, err := t.loadConfig()
	if err == nil {
		return config.AccessCredentials, nil
	}

	tempCred, err := t.oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
	if err != nil {
		return nil, err
	}

	return t.requestClientAuthorization(tempCred)
}

func (t *Twitter) loadConfig() (*config, error) {
	srcName := configFilename()
	src, err := os.Open(srcName)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var loaded *config
	if err := json.NewDecoder(src).Decode(&loaded); err != nil {
		return nil, err
	}

	return loaded, nil
}

type config struct {
	AccessCredentials *oauth.Credentials
}

func (t *Twitter) requestClientAuthorization(tempCred *oauth.Credentials) (*oauth.Credentials, error) {
	url := t.oauthClient.AuthorizationURL(tempCred, nil)
	fmt.Println("open this url: ", url)

	fmt.Print("PIN: ")
	s := bufio.NewScanner(os.Stdin)
	if !s.Scan() {
		return nil, errors.New("failed to scan pin code")
	}
	pin := s.Text()
	token, _, err := t.oauthClient.RequestToken(http.DefaultClient, tempCred, pin)

	return token, err
}

func (t *Twitter) trieve(req *oauthRequest, dest interface{}) error {
	resp, err := t.makeRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return readJSON(resp.Body, dest)
}

func (t *Twitter) makeRequest(req *oauthRequest) (*http.Response, error) {
	params := req.params
	if params == nil {
		params = make(url.Values)
	}

	t.oauthClient.SignParam(req.cred, req.method, req.url, params)

	resp, err := t.doRequest(req.method, req.url, params)
	if err != nil {
		return nil, err
	}
	if http.StatusBadRequest <= resp.StatusCode {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}

func (t *Twitter) doRequest(method, rawURL string, params url.Values) (*http.Response, error) {
	if method != http.MethodGet {
		return http.PostForm(rawURL, params)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	parsed.RawQuery = params.Encode()

	return http.Get(parsed.String())
}

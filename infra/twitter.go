package infra

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/tomocy/mint/domain"
	"github.com/tomocy/mint/infra/twitter"
)

func NewTwitter() *Twitter {
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

func (t *Twitter) FetchHomeTweets() ([]*domain.Tweet, error) {
	cred, err := t.treiveAuthorization()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch home tweets: %s", err)
	}

	var tweets twitter.Tweets
	if err := t.makeRequest(&tweets, cred, http.MethodGet, "https://api.twitter.com/1.1/statuses/home_timeline.json", nil); err != nil {
		return nil, fmt.Errorf("failed to fetch home tweets: %s", err)
	}

	return tweets.Adapt(), nil
}

func (t *Twitter) treiveAuthorization() (*oauth.Credentials, error) {
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

func (t *Twitter) saveConfig(config *config) error {
	destName := configFilename()
	dest, err := os.Open(destName)
	if err != nil {
		return err
	}
	defer dest.Close()

	return json.NewEncoder(dest).Encode(config)
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

func (t *Twitter) makeRequest(dest interface{}, cred *oauth.Credentials, method, rawURL string, params url.Values) error {
	if params == nil {
		params = make(url.Values)
	}

	t.oauthClient.SignParam(cred, method, rawURL, params)

	resp, err := t.doRequest(method, rawURL, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dest)
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

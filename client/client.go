package client

import (
	"fmt"
	"net/http"

	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/config"
)

type WithAuthorizationHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type BacklogClient http.Client

var accessToken string

func SetAccessToken(token string) {
	accessToken = token
}

func (c *BacklogClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{}

	return client.Do(req)
}

func New(config config.BacklogSettings, token string) *backlog.Client {
	httpClient := &BacklogClient{}

	SetAccessToken(token)

	baseUrl := fmt.Sprintf("https://%s.%s", config.Organization, config.BacklogDomain)

	return backlog.New("", baseUrl, backlog.OptionHTTPClient(httpClient))
}

package client

import (
	"fmt"
	"net/http"
)

type WithAuthorizationHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type WithAuthorizationClient http.Client

var accessToken string

func SetAccessToken(token string) {
	accessToken = token
}

func (c *WithAuthorizationClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{}

	return client.Do(req)
}

func Request() {

}

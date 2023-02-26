package main

import (
	"fmt"
	"net/http"

	"github.com/shufo/backlog-cli/cmd"
)

type BacklogSettings struct {
	Project      string `json:"project"`
	Organization string `json:"organization"`
}

type WithAuthorizationHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type WithAuthorizationClient http.Client

var accessToken string

func (c *WithAuthorizationClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{}

	return client.Do(req)

	// if err != nil {
	// 	panic(err)
	// }

	// defer resp.Body.Close()

	// return resp, nil
}

func main() {
	cmd.Execute()

	return
}

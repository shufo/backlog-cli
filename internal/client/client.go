package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/internal/auth"
)

type WithAuthorizationHttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type BacklogClient http.Client

type Error struct {
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"moreInfo"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

var accessToken string

func SetAccessToken(token string) {
	accessToken = token
}

func (c *BacklogClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{}

	// Try to send the request
	resp, _ := client.Do(req)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}

	// parse error response
	var data ErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 400 {
		color.Red(data.Errors[0].Message)
	}

	if resp.StatusCode == 404 && data.Errors[0].Code == 5 {
		color.Red("You are not authorized to access this project.")

		conf, err := config.GetBacklogSetting()

		if err != nil {
			return nil, err
		}

		color.White(fmt.Sprintf("Please verify that you have permission to access %s", conf.Project))
	}

	// If the response status is 401 (Unauthorized) and Code 11, refresh the token and retry
	if resp.StatusCode == 401 && data.Errors[0].Code == 11 {
		conf, err := config.GetBacklogSetting()
		if err != nil {
			return nil, err
		}

		refreshToken, err := config.GetRefreshToken(conf)

		if err != nil {
			return nil, err
		}

		newToken, err := auth.GetAccessTokenFromRefreshToken(conf, refreshToken)

		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", newToken.AccessToken))
		SetAccessToken(newToken.AccessToken)

		hostname := fmt.Sprintf("%s.%s", conf.Organization, conf.BacklogDomain)

		config.WriteConfig(conf.Organization, &config.HostConfig{Hostname: hostname, AccessToken: newToken.AccessToken, RefreshToken: newToken.RefreshToken})

		resp, err = client.Do(req)

		if err != nil {
			return nil, err
		}
	}

	// Return the response after retrying
	return resp, nil

}

func New(conf config.BacklogSettings) *backlog.Client {
	httpClient := &BacklogClient{}

	token, err := config.GetAccessToken(conf)

	if err != nil {
		config.ShowConfigNotFound()
	}

	SetAccessToken(token)

	baseUrl := fmt.Sprintf("https://%s.%s", conf.Organization, conf.BacklogDomain)

	return backlog.New("", baseUrl, backlog.OptionHTTPClient(httpClient))
}

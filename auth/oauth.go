package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/briandowns/spinner"
)

type OauthAuthorizationCodeResponse struct {
	Location string `json:"location"`
	// You can add more fields here if needed
}

var baseUrl = "https://worker-test.shufo.workers.dev"

func GetOauthAuthorizationCode(space string, authCode string) string {
	// Make an HTTP GET request to the API endpoint
	url := fmt.Sprintf("%s/api/v1/oauth?space=%s&auth_code=%s", baseUrl, space, authCode)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Decode the response body into a Response struct
	var response OauthAuthorizationCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	// return location value
	fmt.Println(response.Location)
	return fmt.Sprintf("https://%s.backlog.com%s", space, response.Location)
}

type AuthorizationApprovedResponse struct {
	AuthorizationCode string `json:"code"`
}

func WaitForAuthorizationApprove(space string, authCode string) (string, error) {

	client := http.Client{
		Timeout: time.Second * 10, // Set a timeout for the HTTP request
	}

	startTime := time.Now()

	// Create a new spinner with the default configuration
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// Start the spinner
	s.Start()

	var resp *http.Response
	var data AuthorizationApprovedResponse
	var err error

	for {
		// Make an HTTP GET request to the API endpoint
		endpoint := fmt.Sprintf("%s/api/v1/approve?auth_code=%s", baseUrl, authCode)
		resp, err = client.Get(endpoint)
		if err != nil {
			// If there is an error, wait for 1 second and try again
			time.Sleep(time.Second)
			continue
		}

		// Check the response status code
		if resp.StatusCode == http.StatusOK {
			// If the response code is 200, print a success message and exit the loop
			fmt.Println("Authorization Success!")

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				panic(err)
			}

			break
		}

		// If the response code is not 200, wait for 1 second and try again
		resp.Body.Close()
		time.Sleep(time.Second * 1)

		// Check if the time limit of 10 minutes has been reached
		if time.Since(startTime) > time.Minute*10 {
			fmt.Println("Timeout reached.")
			break
		}
	}

	defer s.Stop()

	if resp.StatusCode == 200 {
		return data.AuthorizationCode, nil
	}

	return "", fmt.Errorf("%s authorization failed", resp.StatusCode)
}

type RequestAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint64 `json:"expires_in"`
}

func GetAccessTokenFromAuthorizationCode(space string, authorizationCode string) RequestAccessTokenResponse {
	var endpoint = fmt.Sprintf("%s/api/v1/token", baseUrl)

	// Set up the POST request parameters for the API
	data := url.Values{}
	data.Set("code", authorizationCode)
	data.Set("space", space)

	// Make the POST request to the OAuth token generation API
	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	var res RequestAccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		panic(err)
	}

	return res
}

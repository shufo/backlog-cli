package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
)

func Login(c *cli.Context) error {
	// Read the backlog.json file
	file, err := os.Open("backlog.json")
	if err != nil {
		return cli.Exit("Error opening backlog.json: "+err.Error(), 1)
	}
	defer file.Close()

	// Decode the JSON settings
	settings := config.BacklogSettings{}
	err = json.NewDecoder(file).Decode(&settings)
	if err != nil {
		return cli.Exit("Error decoding backlog.json: "+err.Error(), 1)
	}

	authCode := util.Genrate7DigitsRandomNumber()

	oauthUrl := GetOauthAuthorizationCode(settings.Organization, authCode)
	// Open the URL in the default web browser
	cmd := exec.Command("open", oauthUrl)
	err = cmd.Run()
	if err != nil {
		return cli.Exit("Error opening URL: "+err.Error(), 1)
	}

	code, err := WaitForAuthorizationApprove(settings.Organization, authCode)
	if err != nil {
		log.Fatalln(err)
	}

	res := GetAccessTokenFromAuthorizationCode(settings.Organization, code)
	fmt.Printf("%s, %s\n", res.AccessToken, res.RefreshToken)
	config.WriteConfig(settings.Organization, &config.HostConfig{AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})

	return nil
}

func RefreshToken() {

}

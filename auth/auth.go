package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
)

func Login(c *cli.Context) error {
	// Read the backlog.json file
	var settings config.BacklogSettings

	if config.ConfigExists() {
		file, err := os.Open("backlog.json")

		if err != nil {
			return cli.Exit("Error opening backlog.json: "+err.Error(), 1)
		}

		defer file.Close()

		// Decode the JSON settings
		err = json.NewDecoder(file).Decode(&settings)
		if err != nil {
			return cli.Exit("Error decoding backlog.json: "+err.Error(), 1)
		}

	} else {

		prompt := promptui.Select{
			Label: "Select your organization's backlog domain",
			Items: []string{"backlog.com", "backlog.jp"},
		}

		_, domain, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil
		}

		fmt.Printf("You choose %q\n", domain)

		validate := func(input string) error {
			if len(input) > 0 {
				return nil
			}

			return errors.New("Invalid string")
		}

		prompt2 := promptui.Prompt{
			Label:    fmt.Sprintf("Input your Organization's name (https://<organization>.%s)", domain),
			Validate: validate,
		}

		org, err := prompt2.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil
		}

		fmt.Printf("You choose %q\n", org)

		settings.BacklogDomain = domain
		settings.Organization = org
	}

	authCode := util.Genrate7DigitsRandomNumber()

	// Create a new spinner with the default configuration
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// Start the spinner
	s.Start()
	defer s.Stop()

	oauthUrl, err := GetOauthAuthorizationUrl(settings, authCode)

	if err != nil {
		s.Stop()
		if err.Error() == "organization not found" {
			fmt.Printf(color.RedString("Organization https://%s.%s/ not found.\n"), settings.Organization, settings.BacklogDomain)
			fmt.Println(color.RedString("Wrong organization name or domain selected."))

			os.Exit(1)
		}

		fmt.Println(err)
		os.Exit(1)
	}

	// Open the URL in the default web browser
	fmt.Printf("Opening %s for Authorization...", oauthUrl)
	cmd := exec.Command("open", oauthUrl)
	err = cmd.Run()
	if err != nil {
		return cli.Exit("Error opening URL: "+err.Error(), 1)
	}

	code, err := WaitForAuthorizationApprove(settings.Organization, authCode)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := GetAccessTokenFromAuthorizationCode(settings.Organization, settings.BacklogDomain, code)

	s.Stop()

	if err != nil {
		log.Fatalln("Authorization Failed.")
	}

	fmt.Println("Authorization Success!")

	hostname := fmt.Sprintf("%s.%s", settings.Organization, settings.BacklogDomain)
	config.WriteConfig(settings.Organization, &config.HostConfig{Hostname: hostname, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})

	config.CreateDefaultConfig(&settings)

	return nil
}

func RefreshToken() {

}

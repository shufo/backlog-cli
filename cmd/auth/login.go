package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/manifoldco/promptui"
	"github.com/shufo/backlog-cli/internal/auth"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/util"
	"github.com/urfave/cli/v3"
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

	authCode, err := util.GenrateUuidV4()

	if err != nil {
		log.Fatalln("Error generating UUID")
	}

	// Create a new spinner with the default configuration
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// Start the spinner
	s.Start()
	defer s.Stop()

	oauthUrl, err := auth.GetOauthAuthorizationUrl(settings, authCode)

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

	util.OpenUrlInBrowser(oauthUrl)

	code, err := auth.WaitForAuthorizationApprove(settings.Organization, authCode, s)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := auth.GetAccessTokenFromAuthorizationCode(settings.Organization, settings.BacklogDomain, code)

	s.Stop()

	if err != nil {
		log.Fatalln("Authorization Failed.", err)
	}

	fmt.Println("Authorization Success!")

	hostname := fmt.Sprintf("%s.%s", settings.Organization, settings.BacklogDomain)
	config.WriteConfig(settings.Organization, &config.HostConfig{Hostname: hostname, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})

	config.CreateDefaultConfig(&settings)

	bl := client.New(settings)

	archived := false
	all := false

	projects, err := bl.GetProjects(&backlog.GetProjectsOptions{Archived: &archived, All: &all})

	if err != nil {
		log.Fatalln(err)
	}

	var projectNames []string

	for _, v := range projects {
		projectNames = append(projectNames, *v.Name)
	}

	prompt := promptui.Select{
		Label: "Select your project",
		Items: projectNames,
		Size:  10,
	}

	_, selectedProject, err := prompt.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		return nil
	}

	fmt.Printf("You choose %q\n", selectedProject)

	var projectKey string

	for _, v := range projects {
		if *v.Name == selectedProject {
			projectKey = *v.ProjectKey
		}
	}

	settings.Project = projectKey

	config.CreateDefaultConfig(&settings)

	return nil
}

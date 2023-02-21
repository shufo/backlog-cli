package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/api"
	"github.com/shufo/backlog-cli/auth"
	"github.com/shufo/backlog-cli/client"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
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
	// result = markdown.Render(string("[aaa](https://example.com)"), 120, 2)
	// fmt.Println(string(result))
	app := &cli.App{
		Name:  "bl",
		Usage: "Open a URL in the default web browser with an embedded parameter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The config file path",
				Value:   "backlog.json",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "issue",
				Usage: "Work with backlog issues",
				Subcommands: []*cli.Command{
					{
						Name:      "view",
						Usage:     "view issue",
						UsageText: "bl issue view <issue_num>",
						ArgsUsage: "<issue_num>",
						Action: func(ctx *cli.Context) error {
							conf, err := config.GetBacklogSetting()

							if err != nil {
								config.ShowConfigNotFound()
								os.Exit(1)
							}

							if ctx.Args().Len() == 0 {
								cli.ShowSubcommandHelpAndExit(ctx, 1)
							}

							id := ctx.Args().First()
							token, err := config.GetAccessToken(conf)

							if err != nil {
								config.ShowConfigNotFound()
							}

							bl := client.New(conf, token)
							api.GetIssue(bl, conf, id)

							return nil
						},
					},
				},
			},
			{
				Name:  "auth",
				Usage: "authentication",
				Subcommands: []*cli.Command{
					{
						Name:      "login",
						Usage:     "Login to backlog organization.\nYou can find organization name at your backlog url https://<organization>.backlog.com/",
						UsageText: "bl auth login <organization>",
						ArgsUsage: "<organization>",
						Action: func(ctx *cli.Context) error {
							auth.Login(ctx)

							return nil
						},
					},
				},
			},
		},
		Action: func(c *cli.Context) error {
			_, err := config.ReadConfig()

			// Check if a parameter was passed
			param := c.Args().Get(0)
			if param == "" {
				return cli.Exit("Please specify an integer parameter.", 1)
			}

			// Read the backlog.json file
			file, err := os.Open("backlog.json")
			if err != nil {
				return cli.Exit("Error opening backlog.json: "+err.Error(), 1)
			}
			defer file.Close()

			// Decode the JSON settings
			settings := BacklogSettings{}
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

			// Construct the URL with the project and parameter embedded
			baseURL := fmt.Sprintf("https://%s.backlog.com/view/", settings.Organization)
			u, err := url.Parse(baseURL)
			if err != nil {
				return cli.Exit("Error parsing base URL: "+err.Error(), 1)
			}

			code, err := WaitForAuthorizationApprove(settings.Organization, authCode)
			if err != nil {
				log.Fatalln(err)
			}

			res := GetAccessTokenFromAuthorizationCode(settings.Organization, code)
			fmt.Printf("%s, %s\n", res.AccessToken, res.RefreshToken)
			accessToken = res.AccessToken
			config.WriteConfig(settings.Organization, &config.HostConfig{AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})

			// config.WriteConfig(settings.Organization, &config.HostConfig{})

			u.Path += settings.Project + "-" + param

			httpClient := &WithAuthorizationClient{}
			bl := backlog.New("", "https://vitgear.backlog.com", backlog.OptionHTTPClient(httpClient))
			// bl := backlog.New("F41Wp4LzLQNmpbhT1xEa84d2Its5KI6ZU1XrsYASSIKrBzUkyLKBYseUr6KNwTzR", "https://vitgear.backlog.com")
			// bl := backlog.New("wpXEIIxJNX4089zt0tyOmreWdo4qVgktb8xmdKPF2mirLmonoXQCrgAzeSffwSr6", "https://vitgear.backlog.com")

			issue, err := bl.GetIssue(fmt.Sprintf("%s-%s", settings.Project, param))
			if err != nil {
				log.Fatalf("%s\n", err)
			}

			if issue == nil {
				log.Fatalln("issue is nil")
			}
			cyan := color.New(color.FgHiCyan)
			reset := color.New(color.Reset)
			cyan.Printf("%s", *issue.Summary)
			reset.Print("")
			cyan.Printf("%s: ", "担当者")
			fmt.Printf("%s\n", *issue.Assignee.Name)
			fmt.Println("")

			// result := markdown.Render(string(*issue.Description), 120, 2)
			// fmt.Println(string(result))
			out, err := glamour.Render(*issue.Description, "dark")
			fmt.Print(out)
			// fullURL := u.String()

			// // Open the URL in the default web browser
			// cmd := exec.Command("open", fullURL)
			// err = cmd.Run()
			// if err != nil {
			// 	return cli.Exit("Error opening URL: "+err.Error(), 1)
			// }

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

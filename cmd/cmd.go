package cmd

import (
	"fmt"
	"os"

	"github.com/shufo/backlog-cli/cmd/auth"
	"github.com/shufo/backlog-cli/cmd/issue"
	"github.com/urfave/cli/v2"
)

func Execute() {
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
						UsageText: "bl issue view <issue_number>",
						ArgsUsage: "<issue_num>",
						Action:    func(ctx *cli.Context) error { return issue.View(ctx) },
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "web",
								Aliases: []string{"w"},
								Usage:   "Open an issue in the browser",
							},
							&cli.BoolFlag{
								Name:    "comments",
								Aliases: []string{"c"},
								Usage:   "View issue comments",
							},
						},
					},
					{
						Name:      "list",
						Usage:     "view issues",
						UsageText: "bl issue list [options]",
						Action:    func(ctx *cli.Context) error { return issue.List(ctx) },
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "assignee",
								Aliases: []string{"a"},
								Usage:   "Filter by assignee",
							},
							&cli.BoolFlag{
								Name:    "me",
								Aliases: []string{"m"},
								Usage:   "Show issues assigned to me",
							},
							&cli.Uint64Flag{
								Name:    "limit",
								Aliases: []string{"L"},
								Usage:   "Maximum number of issues to fetch",
								Value:   15,
							},
							&cli.BoolFlag{
								Name:    "web",
								Aliases: []string{"w"},
								Usage:   "Open an issues in the browser",
							},
						},
					},
					{
						Name:      "status",
						Usage:     "Show status of relevant issues",
						UsageText: "bl issue status [options]",
						Action:    func(ctx *cli.Context) error { return issue.Status(ctx) },
					},
					{
						Name:      "create",
						Usage:     "Create an issue on Backlog.",
						UsageText: "bl issue create [options]",
						Action:    func(ctx *cli.Context) error { return issue.Create(ctx) },
					},
					{
						Name:      "edit",
						Usage:     "Edit an issue on Backlog.",
						UsageText: "bl issue edit [options]",
						Action:    func(ctx *cli.Context) error { return issue.Edit(ctx) },
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
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

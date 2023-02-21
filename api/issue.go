package api

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/auth"
	"github.com/shufo/backlog-cli/config"
)

func GetIssue(bl *backlog.Client, setting config.BacklogSettings, id string) {
	path, err := config.FindConfig("backlog.json")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(path)

	issue, err := bl.GetIssue(fmt.Sprintf("%s-%s", setting.Project, id))

	if err != nil {
		if strings.Contains(err.Error(), "code:11") {
			conf, err := config.GetBacklogSetting()
			if err != nil {
				log.Fatalln(err)
			}
			refreshToken, err := config.GetRefreshToken(conf)

			if err != nil {
				log.Fatalln(err)
			}

			res, err := auth.GetAccessTokenFromRefreshToken(conf, refreshToken)

			if err != nil {
				log.Fatalln("Authorization Failed.")
			}

			hostname := fmt.Sprintf("%s.%s", conf.Organization, conf.BacklogDomain)
			config.WriteConfig(conf.Organization, &config.HostConfig{Hostname: hostname, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken})

			os.Exit(1)
		}
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

	out, err := glamour.Render(*issue.Description, "dark")
	fmt.Print(out)
}

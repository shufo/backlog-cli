package issue

import (
	"fmt"
	"log"
	"os"

	"github.com/shufo/backlog-cli/api"
	"github.com/shufo/backlog-cli/client"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/internal/printer"
	"github.com/urfave/cli/v2"
)

func List(ctx *cli.Context) error {
	conf, err := config.GetBacklogSetting()

	if err != nil {
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	bl := client.New(conf)

	issues, err := api.GetIssueList(bl, conf, ctx)

	if err != nil {
		log.Fatalln(err)
	}

	// print result
	result := printer.PrintIssues(&printer.PrintIssuesParams{Issues: issues})
	fmt.Println(result)

	return nil
}

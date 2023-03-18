package issue

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shufo/backlog-cli/internal/api"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/printer"
	"github.com/urfave/cli/v2"
)

func List(ctx *cli.Context) error {
	// start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	defer s.Stop()
	conf, err := config.GetBacklogSetting()

	if err != nil {
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	bl := client.New(conf)

	issues, err := api.GetIssueList(bl, conf, ctx)

	if err != nil {
		log.Fatalln(err)
		s.Stop()
	}

	s.Stop()

	// print result
	result := printer.PrintIssues(&printer.PrintIssuesParams{Issues: issues})
	fmt.Println(result)

	return nil
}

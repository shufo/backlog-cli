package issue

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shufo/backlog-cli/internal/api"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/printer"
	"github.com/shufo/backlog-cli/internal/util"
	"github.com/urfave/cli/v2"
)

func List(ctx *cli.Context) error {
	// start spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	defer s.Stop()
	conf, err := config.GetBacklogSetting()

	if err != nil {
		s.Stop()
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	bl := client.New(conf)

	if ctx.Bool("web") || util.HasFlag(ctx, "-w", "--web") {
		if err != nil {
			log.Fatalln(err)
		}

		s.Stop()

		url := fmt.Sprintf("https://%s.%s/find/%s?allOver=false&offset=0&order=false&simpleSearch=true&sort=UPDATED", conf.Organization, conf.BacklogDomain, conf.Project)

		cmd := exec.Command("open", url)

		err = cmd.Run()

		if err != nil {
			return cli.Exit("Error opening URL: "+err.Error(), 1)
		}

		return nil
	}

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

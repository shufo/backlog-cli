package issue

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/api"
	"github.com/shufo/backlog-cli/client"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/internal/printer"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
)

func View(ctx *cli.Context) error {
	// start the spinner
	s := spinner.New(spinner.CharSets[14], 500*time.Millisecond)
	s.Start()

	defer s.Stop()

	conf, err := config.GetBacklogSetting()

	if err != nil {
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	if ctx.Args().Len() == 0 {
		cli.ShowSubcommandHelpAndExit(ctx, 1)
	}

	id := ctx.Args().First()

	if ctx.Bool("web") || util.HasFlag(ctx, "-w", "--web") {
		url := fmt.Sprintf("https://%s.%s/view/%s-%s", conf.Organization, conf.BacklogDomain, conf.Project, id)
		cmd := exec.Command("open", url)
		err = cmd.Run()
		if err != nil {
			return cli.Exit("Error opening URL: "+err.Error(), 1)
		}

		return nil
	}

	bl := client.New(conf)

	issue, err := api.GetIssue(bl, conf, id)

	if err != nil {
		s.Stop()
		log.Fatalln(err)
	}

	if issue == nil {
		s.Stop()
		log.Fatalln("issue is nil")
	}

	s.Stop()

	// print result
	printer.PrintIssue(&printer.PrintIssueParams{Issue: issue, Conf: conf, Id: id})

	// comments
	if ctx.Bool("comments") || util.HasFlag(ctx, "-c", "--comments") {
		printComments(&PrintCommentOption{bl, conf, id})
	}

	return nil
}

type PrintCommentOption struct {
	bl   *backlog.Client
	conf config.BacklogSettings
	id   string
}

func printComments(opts *PrintCommentOption) {
	comments, err := api.GetIssueComments(opts.bl, opts.conf, opts.id)

	if err != nil {
		log.Fatalln(err)
	}

	for _, c := range comments {
		if c.Content == nil {
			continue
		}

		printer.PrintIssueComments(&printer.PrintIssueCommentsParams{Comment: c})

	}
}

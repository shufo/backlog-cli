package issue

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/ui"
	"github.com/shufo/backlog-cli/internal/util"
	"github.com/urfave/cli/v3"
)

func Comment(ctx *cli.Context) error {
	conf, err := config.GetBacklogSetting()

	if err != nil {
		log.Fatalln(err)
	}

	bl := client.New(conf)

	if ctx.Args().Len() == 0 {
		fmt.Println(color.RedString("issue id required"))
		os.Exit(1)
	}

	issueId := ctx.Args().First()

	var body string

	body, err = util.GetInputByEditor(&util.GetInputByEditorParam{Current: ""})

	if err != nil {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	issueKey := fmt.Sprintf("%s-%s", conf.Project, issueId)

	submit := ui.Select("What's next?", []string{"Submit", "Cancel"})

	if submit == "Cancel" {
		fmt.Println("Discarded.")
		os.Exit(1)
	}

	if body == "" {
		fmt.Println("There is no change for the issue.")
		os.Exit(1)
	}

	result, err := bl.CreateIssueComment(issueKey, &backlog.CreateIssueCommentInput{Content: &body})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("https://%s.%s/view/%s-%s#comment-%d\n", conf.Organization, conf.BacklogDomain, conf.Project, issueId, *result.ID)

	return nil
}

package issue

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/shufo/backlog-cli/internal/api"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/printer"
	"github.com/urfave/cli/v2"
)

func Status(ctx *cli.Context) error {
	conf, err := config.GetBacklogSetting()

	if err != nil {
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	bl := client.New(conf)

	project, err := bl.GetProject(conf.Project)

	if err != nil {
		log.Fatalln(err)
	}

	user, err := bl.GetUserMySelf()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Relevant issues in %s\n\n", *project.Name)

	issuesAssignedToMe, err := api.GetIssueListAssignedToMe(bl, project, user)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", color.HiBlueString("Issues assigned to you"))
	if len(issuesAssignedToMe) > 0 {
		result := printer.PrintIssues(&printer.PrintIssuesParams{Issues: issuesAssignedToMe})
		fmt.Println(printer.IndentString(result, 2))
	} else {
		fmt.Println(printer.IndentString("No issues assigned to you", 2))
	}

	issuesOpenedByMe, err := api.GetIssueListOpenedByMe(bl, project, user)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", color.HiBlueString("Issues opened by you"))
	if len(issuesOpenedByMe) > 0 {
		result := printer.PrintIssues(&printer.PrintIssuesParams{Issues: issuesOpenedByMe})
		fmt.Println(printer.IndentString(result, 2))
	} else {
		fmt.Println(printer.IndentString("No issues opened by you", 2))
	}

	return nil
}

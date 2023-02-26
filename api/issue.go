package api

import (
	"fmt"
	"log"

	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/config"
	"github.com/urfave/cli/v2"
)

func GetIssue(bl *backlog.Client, setting config.BacklogSettings, id string) (*backlog.Issue, error) {
	return bl.GetIssue(fmt.Sprintf("%s-%s", setting.Project, id))
}

func GetIssueComments(bl *backlog.Client, setting config.BacklogSettings, id string) ([]*backlog.IssueComment, error) {
	options := &backlog.GetIssueCommentsOptions{}

	return bl.GetIssueComments(fmt.Sprintf("%s-%s", setting.Project, id), options)
}

func GetIssueList(bl *backlog.Client, setting config.BacklogSettings, ctx *cli.Context) ([]*backlog.Issue, error) {
	project, err := bl.GetProject(setting.Project)

	if err != nil {
		log.Fatalln(err)
	}

	options := &backlog.GetIssuesOptions{
		Sort: backlog.SortUpdated,
	}

	options.ProjectIDs = []int{*project.ID}

	count := ctx.Int("limit")
	options.Count = &count

	if ctx.Bool("me") {
		user, err := bl.GetUserMySelf()

		if err != nil {
			log.Fatalln(err)
		}

		options.AssigneeIDs = []int{*user.ID}
	}

	return bl.GetIssues(options)
}

func GetIssueListAssignedToMe(bl *backlog.Client, project *backlog.Project, user *backlog.User) ([]*backlog.Issue, error) {
	options := &backlog.GetIssuesOptions{
		Sort: backlog.SortUpdated,
	}

	options.ProjectIDs = []int{*project.ID}

	count := 20
	options.Count = &count

	options.AssigneeIDs = []int{*user.ID}

	return bl.GetIssues(options)
}

func GetIssueListOpenedByMe(bl *backlog.Client, project *backlog.Project, user *backlog.User) ([]*backlog.Issue, error) {
	options := &backlog.GetIssuesOptions{
		Sort: backlog.SortUpdated,
	}

	options.ProjectIDs = []int{*project.ID}

	count := 20
	options.Count = &count

	options.CreatedUserIDs = []int{*user.ID}

	return bl.GetIssues(options)
}

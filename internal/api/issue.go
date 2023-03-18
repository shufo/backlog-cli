package api

import (
	"fmt"
	"log"
	"os"

	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/internal/config"
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

		return bl.GetIssues(options)
	}

	if ctx.String("status") != "" {
		statuses, err := bl.GetStatuses(setting.Project)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var statusId int

		for _, v := range statuses {
			if *v.Name == ctx.String("status") {
				statusId = *v.ID
			}
		}

		options.StatusIDs = []int{statusId}

		return bl.GetIssues(options)
	}

	if ctx.Bool("completed") {
		statuses, err := bl.GetStatuses(setting.Project)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var statusId int

		for _, v := range statuses {
			if *v.Name == "完了" {
				statusId = *v.ID
			}
		}

		options.StatusIDs = []int{statusId}

		return bl.GetIssues(options)
	}

	if ctx.Bool("not-completed") {
		statuses, err := bl.GetStatuses(setting.Project)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var statusIds []int

		for _, v := range statuses {
			if *v.Name != "完了" {
				statusIds = append(statusIds, *v.ID)
			}
		}

		options.StatusIDs = statusIds

		return bl.GetIssues(options)
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

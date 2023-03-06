package printer

import (
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/config"
	"golang.org/x/crypto/ssh/terminal"
)

type PrintIssueParams struct {
	Issue *backlog.Issue
	Conf  config.BacklogSettings
	Id    string
}

const (
	READY       = 1
	IN_PROGRESS = 2
	PROCESSED   = 3
	CLOSED      = 4
)

func PrintIssue(param *PrintIssueParams) {
	// print summary
	cyan := color.New(color.FgHiCyan)
	reset := color.New(color.Reset)
	cyan.Printf("%s ", *param.Issue.Summary)
	reset.Print("")

	// print issue id
	fmt.Printf("%s-%s\n", param.Conf.Project, param.Id)

	// print status, created by, created at
	fmt.Printf(
		"%s %s %s opened about %s\n",
		cyan.Sprintf("状態:"),
		color.New(statusColor(*param.Issue.Status.ID)).Sprint(*param.Issue.Status.Name),
		*param.Issue.CreatedUser.Name,
		TimeDiffString(*param.Issue.Created),
	)

	// print assignee, creator
	var assignee string

	if param.Issue.Assignee != nil {
		assignee = *param.Issue.Assignee.Name
	} else {
		assignee = "-"
	}

	fmt.Printf("%s %s\n", color.CyanString("担当者:"), assignee)

	// print issue type
	fmt.Printf("%s %s\n", color.CyanString("種別:"), *param.Issue.IssueType.Name)

	// print issue priority
	fmt.Printf("%s %s\n", color.CyanString("優先度:"), *param.Issue.Priority.Name)

	// print start date, due date
	if param.Issue.StartDate != nil {
		parsed, err := time.Parse("2006-01-02T00:00:00Z", *param.Issue.StartDate)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s %s\n", color.CyanString("開始日:"), parsed.Format("2006/01/02"))
	}

	if param.Issue.DueDate != nil {
		parsed, err := time.Parse("2006-01-02T00:00:00Z", *param.Issue.DueDate)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s %s\n", color.CyanString("期限日:"), parsed.Format("2006/01/02"))
	}

	fmt.Println("")

	// print description
	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
	)
	out, err := r.Render(*param.Issue.Description)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(out)
}

func statusColor(statusId int) color.Attribute {
	switch statusId {
	case READY:
		return color.FgRed
	case IN_PROGRESS:
		return color.FgBlue
	case PROCESSED:
		return color.FgHiGreen
	case CLOSED:
		return color.FgGreen
	default:
		return color.FgWhite
	}

}

type PrintIssueCommentsParams struct {
	Comment *backlog.IssueComment
}

func PrintIssueComments(param *PrintIssueCommentsParams) {
	fmt.Printf(
		"%s comments %s\n",
		color.HiCyanString(*param.Comment.CreatedUser.Name),
		color.HiBlueString(TimeDiffString(*param.Comment.Created)),
	)

	out, err := glamour.Render(*param.Comment.Content, "dark")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(out)
}

type PrintIssuesParams struct {
	Issues  []*backlog.Issue
	Options table.Options
}

func PrintIssues(param *PrintIssuesParams) string {

	var data [][]string

	for _, v := range param.Issues {
		row := []string{
			*v.IssueKey,
			*v.Summary,
			fmt.Sprintf("about %s", TimeDiffString(*v.Created)),
		}

		data = append(data, row)
	}

	return printTable(data, param.Options)
}

var issueListTableColumnConfig = []table.ColumnConfig{
	{
		Number: 1, Align: text.AlignLeft,
		WidthMin: 12,
		Colors:   text.Colors{text.Color(color.FgHiCyan)},
	},
	{
		Number: 2, Align: text.AlignLeft,
		WidthMax: 120,
		Transformer: func(val interface{}) string {
			chars := utf8.RuneCountInString(val.(string))
			maxChars := 40

			width, err := getTerminalWidth()

			if err != nil {
				log.Fatalln(err)
			}

			switch true {
			case width > 150:
				maxChars = 80
			case width > 100:
				maxChars = 35
			case width > 50:
				maxChars = 20
			}

			if chars > maxChars {
				return truncateString(val.(string), maxChars)
			}

			return val.(string)
		},
	},
	{
		Number: 3, Align: text.AlignLeft,
		WidthMin: 12,
		Colors:   text.Colors{text.Color(color.FgHiBlack)},
	},
}

func printTable(data [][]string, options table.Options) string {
	t := table.NewWriter()
	// t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.Style{Options: table.OptionsNoBordersAndSeparators, Box: table.BoxStyle{PaddingRight: "  "}})

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	for _, row := range data {
		var r table.Row
		for _, v := range row {
			r = append(r, v)
		}
		t.AppendRow(r, rowConfigAutoMerge)

	}

	t.SetColumnConfigs(issueListTableColumnConfig)

	return t.Render()
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

func getTerminalWidth() (int, error) {
	width, _, err := terminal.GetSize(0)
	if err != nil {
		return 0, err
	}
	return width, nil
}

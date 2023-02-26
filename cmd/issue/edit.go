package issue

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/shufo/backlog-cli/client"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/internal/ui"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
)

type model struct {
	options  []string
	selected map[int]bool
	cursor   int
}

type optionMsg int

const (
	Up optionMsg = iota
	Down
	Toggle
	Select
)

type CtrlCMsg struct{}

func Edit(ctx *cli.Context) error {
	conf, err := config.GetBacklogSetting()

	if err != nil {
		log.Fatalln(err)
	}

	bl := client.New(conf)

	fmt.Println(color.HiGreenString("?") + color.HiBlueString(" What would you like to edit?") + " [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")

	var options []string = []string{
		"Issue Type",
		"Summary",
		"Description",
		"Assignee",
	}

	licence, err := bl.GetLicence()

	if err != nil {
		log.Fatalln(err)
	}

	if *licence.Gantt {
		options = append(options, "Start Date")
	}

	options = append(options, "Due Date")

	m := model{
		options:  options,
		selected: make(map[int]bool),
	}

	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	options = m.getSelectedOptions()

	var params *backlog.UpdateIssueInput = &backlog.UpdateIssueInput{}

	issueId := ctx.Args().First()
	issue, err := bl.GetIssue(fmt.Sprintf("%s-%s", conf.Project, issueId))

	if err != nil {
		log.Fatalln(err)
	}

	if util.ContainsString(options, "Issue Type") {
		issueTypeId := getIssueTypeInput(bl, conf)
		params.IssueTypeID = issueTypeId
	}

	if util.ContainsString(options, "Summary") {
		summary := getSummaryInput(&getSummaryInputParam{currentValue: *issue.Summary})
		params.Summary = &summary
	}

	if util.ContainsString(options, "Description") {
		description := getDescriptionInput(&getDescriptionInputParam{currentValue: *issue.Description})
		params.Description = &description
	}

	if util.ContainsString(options, "Assignee") {
		assigneeId := getAssigneeInput(bl, conf)
		if assigneeId > 0 {
			params.AssigneeID = assigneeId
		}
	}

	if util.ContainsString(options, "Start Date") {
		startDate := getStartDateInput()
		params.StartDate = &startDate
	}

	if util.ContainsString(options, "Due Date") {
		dueDate := getDueDateInput()
		params.DueDate = &dueDate
	}

	submit := ui.Select("What's next?", []string{"Submit", "Cancel"})

	if submit == "Cancel" {
		fmt.Println("Discarded.")
		os.Exit(1)
	}

	updatedIssue, err := bl.UpdateIssue(fmt.Sprintf("%s-%s", conf.Project, issueId), params)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("issue %s updated!\n", *updatedIssue.IssueKey)

	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "right":
			for i, _ := range m.options {
				m.selected[i] = true
			}
		case "left":
			for i, _ := range m.options {
				m.selected[i] = false
			}
		case " ":
			if m.selected[m.cursor] {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}
		case "enter":
			return m, tea.Quit
		case "ctrl+c":
			fmt.Println("Canceled")
			os.Exit(1)
		}

	case optionMsg:
		switch msg {
		case Up:
			if m.cursor > 0 {
				m.cursor--
			}
		case Down:
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case Toggle:
			if m.selected[m.cursor] {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}
		case Select:
			return m, tea.Quit
		}
	case CtrlCMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for i, option := range m.options {
		if m.selected[i] {
			if i == m.cursor {
				s += fmt.Sprintf("%s %s %s\n", color.HiCyanString(">"), color.HiGreenString("[x]"), option)
			} else {
				s += fmt.Sprintf("  %s %s\n", color.HiGreenString("[x]"), option)
			}
		} else {
			if i == m.cursor {
				s += fmt.Sprintf("%s %s %s\n", color.HiCyanString(">"), color.HiBlueString("[ ]"), option)
			} else {
				s += "  [ ] " + option + "\n"
			}
		}
	}

	return "\n" + s + "\n"
}

func (m model) getSelectedOptions() []string {
	var selected []string
	for i, isSelected := range m.selected {
		if isSelected {
			selected = append(selected, m.options[i])
		}
	}
	return selected
}

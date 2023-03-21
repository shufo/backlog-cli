package issue

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/manifoldco/promptui"
	"github.com/shufo/backlog-cli/internal/client"
	"github.com/shufo/backlog-cli/internal/config"
	"github.com/shufo/backlog-cli/internal/ui"
	"github.com/shufo/backlog-cli/internal/util"
	"github.com/urfave/cli/v3"
)

type model struct {
	options  []string
	selected map[int]bool
	canceled bool
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
	// start the spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	defer s.Stop()

	conf, err := config.GetBacklogSetting()

	if err != nil {
		log.Fatalln(err)
	}

	bl := client.New(conf)

	var options []string = []string{
		"Issue Type",
		"Summary",
		"Description",
		"Assignee",
		"Status",
		"Priority",
	}

	licence, err := bl.GetLicence()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *licence.Gantt {
		options = append(options, "Start Date")
	}

	options = append(options, "Due Date")

	if ctx.Args().Len() == 0 {
		s.Stop()
		style := lipgloss.NewStyle().
			SetString("Not enough arguments (missing \"Issue id\").\n $ ba issue edit <issue id>").
			Bold(true).
			Foreground(lipgloss.Color("#EFA554")).
			Background(lipgloss.Color("#E356A7")).
			Padding(1)
		fmt.Println(style)
		os.Exit(1)
	}

	issueId := ctx.Args().First()

	issue, err := bl.GetIssue(fmt.Sprintf("%s-%s", conf.Project, issueId))

	if err != nil {
		log.Fatalln(err)
	}

	// stop spinner
	s.Stop()

	fmt.Println(color.HiGreenString("?") + color.HiBlueString(" What would you like to edit?") + " [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]")

	p := tea.NewProgram(model{
		options:  options,
		selected: make(map[int]bool),
	})

	m, err := p.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var selected []string
	if m, ok := m.(model); ok {
		if m.canceled {
			fmt.Println("Canceled")
			os.Exit(1)
		}

		selected = m.getSelectedOptions()
	}

	var params *backlog.UpdateIssueInput = &backlog.UpdateIssueInput{}

	var changed bool

	if util.ContainsString(selected, "Issue Type") {
		issueTypeId := getIssueTypeInput(&getIssueTypeInputParam{bl: bl, conf: conf, current: *issue.IssueType.Name})
		params.IssueTypeID = issueTypeId

		if *issue.IssueType.ID != *issueTypeId {
			changed = true
		}
	}

	if util.ContainsString(selected, "Summary") {
		summary := getSummaryInput(&getSummaryInputParam{current: *issue.Summary})
		params.Summary = &summary

		if *issue.Summary != summary {
			changed = true
		}
	}

	if util.ContainsString(selected, "Description") {
		description, err := getDescriptionInput(&getDescriptionInputParam{current: *issue.Description})

		// skipped
		if err != nil {
			description = *issue.Description
		}

		// if something input
		if description != "" {
			params.Description = &description
		}

		// check if description was changed
		if *issue.Description != description {
			changed = true
		}
	}

	if util.ContainsString(selected, "Assignee") {
		var assigneeId int

		if issue.Assignee != nil {
			assigneeId = getAssigneeInput(&getAssigneeInputParam{bl: bl, conf: conf, current: *issue.Assignee.Name})
		} else {
			assigneeId = getAssigneeInput(&getAssigneeInputParam{bl: bl, conf: conf})
		}

		if assigneeId > 0 {
			params.AssigneeID = assigneeId
		}

		// if the assignee is changed
		if issue.Assignee != nil {
			// if the enter key pressed (skipped)
			if assigneeId == -1 {
				changed = false
			} else if *issue.Assignee.ID != assigneeId {
				// assignee changed
				changed = true
			}
		} else {
			if assigneeId > 0 {
				changed = true
			}
		}

	}

	if util.ContainsString(selected, "Status") {
		statusId := getStatusInput(&getStatusInputParam{bl: bl, conf: conf, current: *issue.Status.Name})
		params.StatusID = &statusId

		if err != nil {
			log.Fatalln(err)
		}

		if *issue.Status.ID != statusId {
			changed = true
		}

	}

	if util.ContainsString(selected, "Priority") {
		priorityId := getPriorityInput(bl)
		params.PriorityID = priorityId

		if err != nil {
			log.Fatalln(err)
		}

		if *issue.Priority.ID != *priorityId {
			changed = true
		}

	}

	if util.ContainsString(selected, "Start Date") {
		var startDate string

		if issue.StartDate != nil {
			startDate = getStartDateInput(&getStartDateInputParam{current: *issue.StartDate})
		} else {
			startDate = getStartDateInput(&getStartDateInputParam{})
		}

		if startDate != "" {
			params.DueDate = &startDate
		}

		// if the start date is changed
		if issue.StartDate != nil {
			parsed, err := time.Parse("2006-01-02T00:00:00Z", *issue.StartDate)

			if err != nil {
				log.Fatalln(err)
			}

			if parsed.Format("2006-01-02") != startDate {
				changed = true
			}

		} else {
			if startDate != "" {
				changed = true
			}
		}
	}

	if util.ContainsString(selected, "Due Date") {
		var dueDate string

		if issue.DueDate != nil {
			dueDate = getDueDateInput(&getDueDateInputParam{current: *issue.DueDate})
		} else {
			dueDate = getDueDateInput(&getDueDateInputParam{})
		}

		if dueDate != "" {
			params.DueDate = &dueDate
		}

		// if the due date is changed
		if issue.DueDate != nil {
			parsed, err := time.Parse("2006-01-02T00:00:00Z", *issue.DueDate)

			if err != nil {
				log.Fatalln(err)
			}

			if parsed.Format("2006-01-02") != dueDate {
				changed = true
			}

		} else {
			if dueDate != "" {
				changed = true
			}
		}
	}

	submit := ui.Select("What's next?", []string{"Submit", "Cancel"})

	if submit == "Cancel" {
		fmt.Println("Discarded.")
		os.Exit(1)
	}

	if !changed {
		fmt.Println("There is no change for the issue.")
		os.Exit(1)
	}

	updatedIssue, err := bl.UpdateIssue(fmt.Sprintf("%s-%s", conf.Project, issueId), params)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("issue %s updated!\n", *updatedIssue.IssueKey)
	fmt.Printf("https://%s.%s/view/%s\n", conf.Organization, conf.BacklogDomain, *updatedIssue.IssueKey)

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
			m.canceled = true
			return m, tea.Quit
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

type getStatusInputParam struct {
	bl      *backlog.Client
	conf    config.BacklogSettings
	current string
}

func getStatusInput(param *getStatusInputParam) int {
	// get issue types
	statuses, err := param.bl.GetStatuses(param.conf.Project)

	if err != nil {
		log.Fatalln(err)
	}

	var items []string

	for _, v := range statuses {
		items = append(items, *v.Name)
	}

	var cursorPos int

	if param.current != "" {

		for i, v := range statuses {
			if *v.Name == param.current {
				cursorPos = i
			}
		}
	}

	var label string

	if param.current != "" {
		label = fmt.Sprintf("Select status (%s)", param.current)
	} else {
		label = "Select status"
	}

	promptStatus := promptui.Select{
		Label: label,
		Items: items,
		Size:  10,
	}

	_, selected, err := promptStatus.RunCursorAt(cursorPos, 0)

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("You choose status as %q\n", selected)

	var selectedId int

	for _, v := range statuses {
		if *v.Name == selected {
			selectedId = *v.ID
		}
	}

	return selectedId
}

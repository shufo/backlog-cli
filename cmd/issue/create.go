package issue

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/kenzo0107/backlog"
	"github.com/manifoldco/promptui"
	"github.com/shufo/backlog-cli/client"
	"github.com/shufo/backlog-cli/config"
	"github.com/shufo/backlog-cli/internal/ui"
	"github.com/shufo/backlog-cli/util"
	"github.com/urfave/cli/v2"
)

// Create is a function that creates a new issue in a backlog management tool, like Jira or Asana. It takes a command-line context object as an input parameter.
func Create(ctx *cli.Context) error {
	// Get configuration settings for the backlog tool
	conf, err := config.GetBacklogSetting()

	// If there is an error getting the settings, print an error message and exit the program
	if err != nil {
		config.ShowConfigNotFound()
		os.Exit(1)
	}

	// Create a new client object to interact with the backlog tool
	bl := client.New(conf)

	// Get the project from the backlog tool
	project, err := bl.GetProject(conf.Project)

	// If there is an error getting the project, log the error and exit the program
	if err != nil {
		log.Fatalln(err)
	}

	// Prompt the user to input the issue type, summary, description, priority, and assignee of the new issue, using various get functions
	issueTypeId := getIssueTypeInput(bl, conf)
	summary := getSummaryInput(&getSummaryInputParam{})
	description := getDescriptionInput(&getDescriptionInputParam{})
	priorityId := getPriorityInput(bl)
	assigneeId := getAssigneeInput(bl, conf)

	licence, err := bl.GetLicence()

	if err != nil {
		log.Fatalln(err)
	}

	var startDate string

	if *licence.Gantt {
		startDate = getStartDateInput()
	}

	dueDate := getDueDateInput()

	submit := ui.Select("What's next?", []string{"Submit", "Cancel"})

	if submit == "Cancel" {
		fmt.Println("Discarded.")
		os.Exit(1)
	}

	// Create a backlog.CreateIssueInput object that contains the necessary parameters to create the issue
	params := &backlog.CreateIssueInput{
		ProjectID:   project.ID,
		Summary:     &summary,
		Description: &description,
		IssueTypeID: issueTypeId,
		PriorityID:  priorityId,
		StartDate:   &startDate,
		DueDate:     &dueDate,
	}

	// If there is an assignee ID provided, assign the value to the params.AssigneeID field
	if assigneeId > 0 {
		params.AssigneeID = &assigneeId
	}

	// Create the issue using bl.CreateIssue(params) function and log the issue key and summary if the issue is created successfully
	issue, err := bl.CreateIssue(params)
	if err != nil {
		// If there is an error creating the issue, log the error and exit the program
		log.Fatalln(err)
	}

	if issue != nil {
		fmt.Printf("issue %s %s created\n", *issue.IssueKey, summary)
		fmt.Printf("https://%s.%s/view/%s\n", conf.Organization, conf.BacklogDomain, *issue.IssueKey)
	} else {
		fmt.Printf("issue %s created", summary)
	}

	return nil
}

func getIssueTypeInput(bl *backlog.Client, conf config.BacklogSettings) *int {
	// get issue types
	issueTypes, err := bl.GetIssueTypes(conf.Project)

	if err != nil {
		log.Fatalln(err)
	}

	var it []string

	for _, v := range issueTypes {
		it = append(it, *v.Name)
	}

	promptIssueType := promptui.Select{
		Label: "Select issue type",
		Items: it,
		Size:  10,
	}

	_, selectedIssueType, err := promptIssueType.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("You choose %q\n", selectedIssueType)

	var selectedIssueTypeId *int

	for _, v := range issueTypes {
		if *v.Name == selectedIssueType {
			selectedIssueTypeId = v.ID
		}
	}

	return selectedIssueTypeId
}

type getSummaryInputParam struct {
	currentValue string
}

func getSummaryInput(param *getSummaryInputParam) string {
	validate := func(input string) error {
		if len(input) > 0 && utf8.RuneCountInString(input) <= 255 {
			return nil
		}

		if len(input) == 0 {
			return errors.New("summary needs at least 1 characters")
		}

		if utf8.RuneCountInString(input) > 255 {
			return errors.New("summary must be within 255 characters")
		}

		return errors.New("summary must be betwen 1 and 255 characters")
	}

	var label string

	if param.currentValue != "" {
		label = fmt.Sprintf("Summary (%s)", param.currentValue)
	} else {
		label = "Summary"
	}

	promptSummary := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	summary, err := promptSummary.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		os.Exit(1)
	}

	return summary
}

type getDescriptionInputParam struct {
	currentValue string
}

func getDescriptionInput(param *getDescriptionInputParam) string {
	editor := util.DetectEditor()

	fmt.Printf(
		"%s %s\n",
		color.HiGreenString("?"),
		color.BlueString(fmt.Sprintf("Body [(e) to launch %s, enter to skip]", path.Base(editor))),
	)

	char, key, err := waitForKey(&waitForKeyInput{
		keys: []keyboard.Key{
			keyboard.KeyEnter,
			keyboard.KeyCtrlC,
		},
		chars: []rune{
			'e',
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	var editKeyCode rune = 'e' // Replace with the desired key code in rune format (e.g., '1' for the '1' key)

	var description string

	if char == editKeyCode {
		description, err = openEditor(param.currentValue)

		if err != nil {
			log.Fatalln(err)
		}
	}

	if utf8.RuneCountInString(description) > 100_000 {
		fmt.Println(color.RedString("Issue description must be within 100,000 characters."))
		os.Exit(1)
	}

	if key == keyboard.KeyCtrlC {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	return description
}

func getPriorityInput(bl *backlog.Client) *int {
	priorities, err := bl.GetPriorities()

	if err != nil {
		log.Fatalln(err)
	}

	var pr []string

	for _, v := range priorities {
		pr = append(pr, *v.Name)
	}

	promptPriority := promptui.Select{
		Label: "Select priority",
		Items: pr,
		Size:  10,
	}

	_, selectedPriority, err := promptPriority.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		return nil
	}

	var selectedPriorityId *int
	for _, v := range priorities {
		if *v.Name == selectedPriority {
			selectedPriorityId = v.ID
		}
	}

	return selectedPriorityId
}

func getAssigneeInput(bl *backlog.Client, conf config.BacklogSettings) int {
	fmt.Printf(
		"%s %s\n",
		color.HiGreenString("?"),
		color.BlueString("Assignee [(s) to select, (m) to assign self, enter to skip]"),
	)

	targetRunes := []rune{'s', 'm'}

	char, key, err := waitForKey(&waitForKeyInput{
		keys: []keyboard.Key{
			keyboard.KeyEnter,
			keyboard.KeyCtrlC,
		},
		chars: targetRunes,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Skip answer if key is enter
	if key == keyboard.KeyEnter {
		return -1
	}

	// Cancel flow if key is Ctrl+C
	if key == keyboard.KeyCtrlC {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	// Assign self if entered char is `m`
	if char == 'm' {
		me, err := bl.GetUserMySelf()

		if err != nil {
			log.Fatalln(err)
		}

		return *me.ID
	}

	// Select user if entered char is `s`
	excludeGroupMembers := false

	users, err := bl.GetProjectUsers(conf.Project, &backlog.GetProjectUsersOptions{ExcludeGroupMembers: &excludeGroupMembers})

	if err != nil {
		log.Fatalln(err)
	}

	var memberNames []string

	for _, v := range users {
		memberNames = append(memberNames, *v.Name)
	}

	promptPriority := promptui.Select{
		Label: "Select assignee",
		Items: memberNames,
		Size:  10,
	}

	_, selectedMember, err := promptPriority.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		os.Exit(1)
	}

	var selectedMemberId int
	for _, v := range users {
		if *v.Name == selectedMember {
			selectedMemberId = *v.ID
		}
	}

	return selectedMemberId
}

func openEditor(currentValue string) (string, error) {

	var editor string
	editor = os.Getenv("EDITOR")

	if editor == "" {
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		case "darwin":
			editor = "vim"
		case "linux":
			editor = "vim"
		default:
			return "", errors.New("unsupported operating system")
		}

	}

	// Create a temporary file for the user to edit
	tempFile, err := ioutil.TempFile("", "bl_issue_create")
	if err != nil {
		return "", err
	}
	// Write initial value
	tempFile.WriteString(currentValue)

	defer os.Remove(tempFile.Name()) // Clean up the temporary file when done

	// Launch Vim to edit the temporary file
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	// Read the contents of the saved file
	contents, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

type waitForKeyInput struct {
	keys  []keyboard.Key
	chars []rune
}

func waitForKey(input *waitForKeyInput) (rune, keyboard.Key, error) {
	err := keyboard.Open()

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetSingleKey()

		if err != nil {
			panic(err)
		}

		for _, v := range input.keys {
			if key == v {
				return char, key, err
			}
		}

		for _, v := range input.chars {
			if char == v {
				return char, key, err
			}

		}
	}
}

func getStartDateInput() string {
	fmt.Printf(
		"%s %s\n",
		color.HiGreenString("?"),
		color.BlueString("Start Date [(s) to select, enter to skip]"),
	)

	targetRunes := []rune{'s'}

	_, key, err := waitForKey(&waitForKeyInput{
		keys: []keyboard.Key{
			keyboard.KeyEnter,
			keyboard.KeyCtrlC,
		},
		chars: targetRunes,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Skip answer if key is enter
	if key == keyboard.KeyEnter {
		return ""
	}

	// Cancel flow if key is Ctrl+C
	if key == keyboard.KeyCtrlC {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	// if s pressed
	dates := generateDates(30)

	promptPriority := promptui.Select{
		Label: "Select Start Date",
		Items: dates,
		Size:  10,
	}

	_, selectedDate, err := promptPriority.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		return ""
	}

	return selectedDate
}

func getDueDateInput() string {
	fmt.Printf(
		"%s %s\n",
		color.HiGreenString("?"),
		color.BlueString("Due Date [(s) to select, enter to skip]"),
	)

	targetRunes := []rune{'s'}

	_, key, err := waitForKey(&waitForKeyInput{
		keys: []keyboard.Key{
			keyboard.KeyEnter,
			keyboard.KeyCtrlC,
		},
		chars: targetRunes,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Skip answer if key is enter
	if key == keyboard.KeyEnter {
		return ""
	}

	// Cancel flow if key is Ctrl+C
	if key == keyboard.KeyCtrlC {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	// if s pressed
	dates := generateDates(30)

	promptPriority := promptui.Select{
		Label: "Select Due Date",
		Items: dates,
		Size:  10,
	}

	_, selectedDate, err := promptPriority.Run()

	if err != nil {
		fmt.Printf("Canceled %v\n", err)
		return ""
	}

	return selectedDate
}

func generateDates(days int) []string {
	// Create an empty slice to store the generated dates
	dates := []string{}

	// Generate the dates
	for i := 0; i < days; i++ {
		// Calculate the date i days in the future
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")

		// Append the date to the slice
		dates = append(dates, date)
	}

	// return the generated dates
	return dates
}

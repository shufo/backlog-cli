package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

var choices []string

type model struct {
	cursor   int
	choice   string
	canceled bool
}

func Select(title string, chs []string) string {
	fmt.Printf("%s %s\n", color.HiGreenString("?"), color.HiBlueString(title))
	p := tea.NewProgram(model{})
	choices = chs

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok {
		if m.canceled {
			fmt.Println("Canceled")
			os.Exit(1)
		}

		if m.choice != "" {
			return m.choice

		}
	}

	return ""
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.canceled = true
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString(fmt.Sprintf("%s %s", color.HiCyanString(">"), color.HiCyanString(choices[i])))
		} else {
			s.WriteString(fmt.Sprintf("  %s", choices[i]))
		}
		s.WriteString("\n")
	}
	s.WriteString("\n")

	return s.String()
}

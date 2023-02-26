package printer

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/kenzo0107/backlog"
)

func TimeDiffString(t backlog.Timestamp) string {
	now := time.Now()
	diff := now.Sub(t.Time)

	seconds := int(diff.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30

	switch {
	case months > 0:
		return fmt.Sprintf("%d months ago", months)
	case days > 1:
		return fmt.Sprintf("%d days ago", days)
	case days == 1:
		return "yesterday"
	case hours > 1:
		return fmt.Sprintf("%d hours ago", hours)
	case hours == 1:
		return "an hour ago"
	case minutes > 1:
		return fmt.Sprintf("%d minutes ago", minutes)
	case minutes == 1:
		return "a minute ago"
	default:
		return "just now"
	}
}

func IndentString(input string, spaces int) string {
	// Split the input into individual lines
	lines := strings.Split(input, "\n")

	// Create a buffer to store the output
	var buffer bytes.Buffer

	// Iterate over each line and add the specified number of spaces
	for _, line := range lines {
		buffer.WriteString(strings.Repeat(" ", spaces))
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}

	// Return the indented string
	return buffer.String()
}

package util

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"
)

func Genrate7DigitsRandomNumber() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(10000000)

	firstPart := num / 1000  // take the first three digits of the number
	secondPart := num % 1000 // take the last three digits of the number

	return fmt.Sprintf("%04d-%03d", firstPart, secondPart)
}

func HasFlag(ctx *cli.Context, flags ...string) bool {
	for _, v := range flags {
		if ContainsString(ctx.Args().Slice(), v) {
			return true
		}
	}

	return false
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func DetectEditor() string {
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
			return ""
		}

	}

	return editor
}

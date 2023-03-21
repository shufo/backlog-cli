package util

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/urfave/cli/v3"
)

func GenrateUuidV4() (string, error) {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(crand.Reader, uuid)

	if err != nil {
		return "", err
	}

	// Set version (4) and variant (2)
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
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

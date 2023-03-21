package util

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
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

type GetInputByEditorParam struct {
	Current string
}

func GetInputByEditor(param *GetInputByEditorParam) (string, error) {
	editor := DetectEditor()

	fmt.Printf(
		"%s %s\n",
		color.HiGreenString("?"),
		color.BlueString(fmt.Sprintf("Body [(e) to launch %s]", path.Base(editor))),
	)

	char, key, err := waitForKey(&waitForKeyInput{
		keys: []keyboard.Key{
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
		description, err = openEditor(param.Current)

		if err != nil {
			log.Fatalln(err)
		}
	}

	if utf8.RuneCountInString(description) > 100_000 {
		fmt.Println(color.RedString("Input must be within 100,000 characters."))
		os.Exit(1)
	}

	if key == keyboard.KeyCtrlC {
		fmt.Println("Canceled")
		os.Exit(1)
	}

	return description, nil
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
	tempFile, err := ioutil.TempFile("", "tmp_")
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

func OpenUrlInBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer.exe", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		if IsWSL() {
			cmd = exec.Command("cmd.exe", "/C", "start", url)
		} else {
			cmd = exec.Command("xdg-open", url)
		}
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	cmd.Run()
}

func IsWSL() bool {
	_, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop")

	return !os.IsNotExist(err)
}

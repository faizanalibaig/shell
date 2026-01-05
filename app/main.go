package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type builtin int

const (
	echo builtin = iota
	exit
	type_
	pwd
	cd
)

var builtins = map[string]bool{
	echo.String():  true,
	exit.String():  true,
	type_.String(): true,
	pwd.String():   true,
	cd.String():    true,
}

func (b builtin) String() string {
	switch b {
	case echo:
		return "echo"
	case exit:
		return "exit"
	case type_:
		return "type"
	case pwd:
		return "pwd"
	case cd:
		return "cd"
	default:
		return "unknown"
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		prompt, err := ReadFromStdin()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		if len(prompt) == 0 {
			continue
		}

		cmd := prompt[0]
		args := prompt[1:]

		switch cmd {
		case exit.String():
			HandleExit()
		case echo.String():
			HandleEcho(args)
		case type_.String():
			CheckType(args[0])
		case pwd.String():
			GetCurrentDir()
		case cd.String():
			HandleChangeDir(args[0])
		default:
			err := ExecuteCommand(cmd, args...)
			if err != nil {
				fmt.Printf("%v: command not found\n", cmd)
			}
		}
	}
}

func ReadFromStdin() ([]string, error) {
	prompt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nil, err
	}

	command := strings.TrimRight(prompt, "\r\n")
	return parseInput(command), nil
}

func parseInput(input string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '\'':
			inSingleQuote = !inSingleQuote

		case ' ':
			if inSingleQuote {
				current.WriteByte(ch)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}

		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func CheckType(cmd string) {
	if ok := builtins[cmd]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	}

	fullPath, ok := GetFullPath(cmd)
	if ok {
		fmt.Printf("%s is %s\n", cmd, fullPath)
		return
	}

	fmt.Printf("%s: not found\n", cmd)
}

func GetFullPath(cmd string) (string, bool) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}

	return path, true
}

func ExecuteCommand(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command.Run()
}

func GetCurrentDir() {
	dir, _ := os.Getwd()
	fmt.Printf("%s\n", dir)
}

func HandleChangeDir(path string) {
	_, err := os.Stat(path)

	if path == "~" {
		home := os.Getenv("HOME")
		HandleHomeDir(home)
	} else if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
	} else {
		_ = os.Chdir(path)
	}
}

func HandleHomeDir(home string) {
	_ = os.Chdir(home)
}

func HandleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func HandleExit() {
	os.Exit(0)
}

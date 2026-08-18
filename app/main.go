package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
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

var stdin = bufio.NewReader(os.Stdin)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		prompt, err := ReadFromStdin()

		if err == io.EOF {
			os.Exit(0)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		if len(prompt) == 0 {
			continue
		}

		cmd := prompt[0]
		args := prompt[1:]

		var stdout io.Writer = os.Stdout
		var redirectFile *os.File

		rest, fileName, hasRedirect := ParseRedirect(args)
		if hasRedirect {
			file, err := OpenRedirectFile(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", fileName)
				continue
			}

			redirectFile = file
			stdout = file
			args = rest
		}

		switch cmd {
		case exit.String():
			HandleExit()
		case echo.String():
			HandleEcho(stdout, args)
		case type_.String():
			CheckType(stdout, args[0])
		case pwd.String():
			GetCurrentDir(stdout)
		case cd.String():
			HandleChangeDir(args[0])
		default:
			err := ExecuteCommand(stdout, cmd, args...)
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Printf("%v: command not found\n", cmd)
			}
		}

		if redirectFile != nil {
			redirectFile.Close()
		}
	}
}

func ReadFromStdin() ([]string, error) {
	prompt, err := stdin.ReadString('\n')
	if err != nil {
		return nil, err
	}

	command := strings.TrimRight(prompt, "\r\n")
	token, err := shlex.Split(command)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ParseRedirect(args []string) ([]string, string, bool) {
	for i, arg := range args {
		if arg != ">" && arg != "1>" {
			continue
		}

		if i+1 >= len(args) {
			return args, "", false
		}

		return args[:i], args[i+1], true
	}

	return args, "", false
}

func OpenRedirectFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func CheckType(stdout io.Writer, cmd string) {
	if ok := builtins[cmd]; ok {
		fmt.Fprintf(stdout, "%s is a shell builtin\n", cmd)
		return
	}

	fullPath, ok := GetFullPath(cmd)
	if ok {
		fmt.Fprintf(stdout, "%s is %s\n", cmd, fullPath)
		return
	}

	fmt.Fprintf(stdout, "%s: not found\n", cmd)
}

func GetFullPath(cmd string) (string, bool) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}

	return path, true
}

func ExecuteCommand(stdout io.Writer, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Stdout = stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command.Run()
}

func GetCurrentDir(stdout io.Writer) {
	dir, _ := os.Getwd()
	fmt.Fprintf(stdout, "%s\n", dir)
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

func HandleEcho(stdout io.Writer, args []string) {
	fmt.Fprintln(stdout, strings.Join(args, " "))
}

func HandleExit() {
	os.Exit(0)
}

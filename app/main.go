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
				fmt.Printf("%v: command not found \n", cmd)
			}
		}
	}
}

func ReadFromStdin() ([]string, error) {
	prompt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return []string{}, fmt.Errorf("error reading from stdin: %v", err)
	}

	command := strings.TrimRight(prompt, "\r\n")
	return strings.Split(command, " "), nil
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

func GetCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fmt.Printf("%s\n", dir)
	return dir, nil
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
	joined := strings.Join(args, " ")
	fmt.Println(joined)
}

func HandleExit() {
	os.Exit(0)
}

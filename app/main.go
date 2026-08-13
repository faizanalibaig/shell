package main

import (
	"bufio"
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
	cat
	type_
	pwd
	cd
)

var builtins = map[string]bool{
	echo.String():  true,
	cat.String():   true,
	exit.String():  true,
	type_.String(): true,
	pwd.String():   true,
	cd.String():    true,
}

func (b builtin) String() string {
	switch b {
	case echo:
		return "echo"
	case cat:
		return "cat"	
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

		
		switch cmd {
		case exit.String():
			HandleExit()
		case echo.String():
			if len(args) > 2 && (args[len(args) - 2] == ">" || args[len(args) - 2] == "1>") {
				fileName := args[len(args)-1]
				content := strings.Join(args[:len(args)-2], " ")
				RedirectOutputToFile(fileName, content)
				continue
			} else {
				HandleEcho(args)
			}
		case cat.String():
			if len(args) <= 1 {
				ReadContentFromFile(args[0])	
			} else {
				fmt.Fprintln(os.Stderr, "cat: too many arguments")
				os.Exit(1)
			}
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

func RedirectOutputToFile(fileName string, content string) {
	err := os.WriteFile(fileName, []byte(content), 0644)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to file:", err)
		os.Exit(1)
	}
}

func ReadContentFromFile(fileName string) {
    if len(fileName) == 0 {
		fmt.Errorf("No file name provided")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(string(content))
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

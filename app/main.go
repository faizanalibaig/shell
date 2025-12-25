package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type builtin int

const (
	echo builtin = iota
	exit
	type_
)

var builtins = map[string]bool{
	echo.String():  true,
	exit.String():  true,
	type_.String(): true,
}

func (b builtin) String() string {
	switch b {
	case echo:
		return "echo"
	case exit:
		return "exit"
	case type_:
		return "type"
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

		path := os.Getenv("PATH")

		switch cmd {
		case exit.String():
			HandleExit()
		case echo.String():
			HandleEcho(args)
		case type_.String():
			CheckType(args[0], path)
		default:
			fmt.Printf("%v: command not found \n", cmd)
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

func CheckType(cmd string, path string) {
	if ok := builtins[cmd]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	}

	fullPath, ok := GetFullPath(cmd, path)
	if ok {
		fmt.Printf("%s is %s\n", cmd, fullPath)
		return
	}

	fmt.Printf("%s: not found\n", cmd)
}

func GetFullPath(cmd, path string) (string, bool) {
	paths := strings.Split(path, ";")

	for _, p := range paths {
		fullPath := filepath.Join(p, cmd)
		if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
			fmt.Println(p)
			return fullPath, true
		}
	}

	return "", false
}

func HandleEcho(args []string) {
	joined := strings.Join(args, " ")
	fmt.Println(joined)
}

func HandleExit() {
	os.Exit(0)
}

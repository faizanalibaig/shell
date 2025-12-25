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

		switch cmd {
		case exit.String():
			HandleExit()
		case echo.String():
			HandleEcho(args)
		case type_.String():
			CheckType(args[0])
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
	output, err := command.Output()

	if err != nil {
		return err
	}

	fmt.Printf("Program was passed %v args\n Arg #0 (program name): %s", len(args)+1, cmd)
	for index, line := range args {
		fmt.Printf(" Args #%v: %s\n", index+1, line)
	}

	fmt.Printf("Program Signature: %v", string(output))
	return nil
}

func HandleEcho(args []string) {
	joined := strings.Join(args, " ")
	fmt.Println(joined)
}

func HandleExit() {
	os.Exit(0)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	for true {
		fmt.Print("$ ")
		prompt, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command := strings.TrimSpace(prompt)

		if command[:5] == "type " {
			if len(command) > 4 {
				checkType(command[5:])
			}

			continue
		} else if command == "exit" {
			os.Exit(0)
		} else if command[:5] == "echo " {
			if len(command) > 4 {
				echo(command[5:])
			}

			continue
		}

		fmt.Printf("%v: command not found \n", command)
	}
}

func checkType(message string) {
	types := []string{"echo", "exit", "type"}

	if slices.Contains(types, message) {
		fmt.Printf("%v is a shell builtin \n", message)
	} else {
		fmt.Printf("%v: not found \n", message)
	}
}

func echo(message string) {
	fmt.Println(message)
}

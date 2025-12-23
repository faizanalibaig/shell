package main

import (
	"bufio"
	"fmt"
	"os"
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

		if command == "exit" {
			os.Exit(0)
		}

		fmt.Printf("%v: command not found \n", command)
	}
}

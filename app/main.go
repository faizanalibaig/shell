package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Print("$ ")
	prompt, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	fmt.Printf("%v: command not found", strings.TrimSpace(prompt))
}

package main

import "fmt"

func main() {
	var prompt string

	fmt.Print("$ ")
	_, err := fmt.Scan(&prompt)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

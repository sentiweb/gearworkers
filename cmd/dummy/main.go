package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	command := "hello"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if command == "hello" {
		fmt.Println("Hello")
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
	}

	if command == "echo" {
		if len(os.Args) > 2 {
			a := os.Args[2:]
			fmt.Println(strings.Join(a, " "))
		}
	}
	if command == "error" {
		os.Exit(2)
	}
}

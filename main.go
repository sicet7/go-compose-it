package main

import (
	"fmt"
	"go-compose-it/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("command execution retuned error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

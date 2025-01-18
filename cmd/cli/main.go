package main

import (
	"log"

	command "konbini/cli/commands"
)

func main() {
	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

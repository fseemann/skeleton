package main

import (
	"github.com/skeleton/skeleton-cli/add"
	"log"
	"os"
)

var cmd = []string{"add"}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Use on of the following commands: \n", cmd)
	}
	switch args[1] {
	case "add":
		add.Add()
	}
}

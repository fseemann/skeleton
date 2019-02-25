package main

import (
	"github.com/skeleton/skeleton-cli/template"
	"log"
	"os"
)

var cmd = []string{"create"}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Use on of the following commands: \n", cmd)
	}

	remainingArgs := args[2:]
	switch args[1] {
	case "create":
		template.Create(remainingArgs)
	case "build":
		template.Build(remainingArgs)
	default:
		log.Fatalf("usage: %v", cmd)
	}
}

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
	switch args[1] {
	case "create":
		template.Create(args[2:])
	default:
		log.Fatalf("usage: %v", cmd)
	}
}

package main

import (
	"github.com/Kenedy228/blog-aggregator/internal/commands"
	"github.com/Kenedy228/blog-aggregator/internal/handlers"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalf("provide command")
	}

	cmd := commands.NewCommand(args[0], args[1:])
	err := handlers.ResolveCommand(cmd)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

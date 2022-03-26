package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bryanro/repo-stats/pkg/stats"
	"github.com/bryanro/repo-stats/pkg/version"
)

func main() {
	ctx := context.Background()
	log.Printf("starting, git commit %s", version.GitCommit)

	// Parse our input and create our options struct
	options, err := stats.CheckArgs(os.Args[1:])
	if err != nil {
		// Print the error first and then show an example
		// of how to run the command and exit 1
		fmt.Println(err)
		stats.Usage()
		os.Exit(1)
	}

	// Enter into our program, exit on error
	if err := stats.Run(ctx, options); err != nil {
		log.Fatal(err)
	}
}

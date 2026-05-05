package main

import (
	"context"
	"fmt"
	"interFleet/src"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	app := src.App{}

	args := os.Args[1:]
	fmt.Printf("Args: %v\n", args)
	for i := range args {
		if (args[i] == "-h" || args[i] == "--host") && len(args) > i+1 {
			app.HostURLConfig = args[i+1]
		}
		if (args[i] == "-f" || args[i] == "--file") && len(args) > i+1 {
			app.CSVPathConfig = args[i+1]
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

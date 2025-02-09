package main

import (
	"context"
	"log"

	cli "github.com/monochromegane/mcp-bigquery-server"
)

func main() {
	ctx := context.TODO()
	if err := run(ctx); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(ctx context.Context) error {
	c, err := cli.New()
	if err != nil {
		return err
	}

	return c.Run(ctx)
}

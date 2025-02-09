package mcp_bigquery_server

import (
	"context"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Start struct {
		Project  string   `required:"" help:"Project ID"`
		Location string   `default:"asia-northeast1" help:"Location"`
		Dataset  []string `required:"" help:"Allowed datasets"`
	} `cmd:"" help:"Start the MCP BigQuery server"`
}

func New() (*CLI, error) {
	return &CLI{}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	k := kong.Parse(c)

	switch k.Command() {
	case "start":
		return StartServer(ctx, c)
	}
	return nil
}

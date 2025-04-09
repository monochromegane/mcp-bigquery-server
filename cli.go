package mcp_bigquery_server

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Version kong.VersionFlag `help:"Show version"`
	Start   struct {
		Project string   `required:"" help:"Project ID"`
		Dataset []string `required:"" help:"Allowed datasets"`
	} `cmd:"" help:"Start the MCP BigQuery server"`
}

func New() (*CLI, error) {
	return &CLI{}, nil
}

func (c *CLI) Run(ctx context.Context) error {
	k := kong.Parse(c, kong.Vars{
		"version": fmt.Sprintf("%s v%s (rev:%s)", "mcp-bigquery-server", version, revision),
	})

	switch k.Command() {
	case "start":
		return StartServer(ctx, c)
	}
	return nil
}

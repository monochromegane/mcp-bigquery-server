package mcp_bigquery_server

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolName string

const (
	LIST_TABLES      ToolName = "list_tables"
	GET_TABLE_SCHEMA ToolName = "get_table_schema"
)

func StartServer(ctx context.Context, c *CLI) error {
	bs, err := NewBigQueryServer(ctx, c.Start.Project, c.Start.Location, c.Start.Dataset)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := bs.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	return nil
}

func NewBigQueryServer(ctx context.Context, project, location string, datasets []string) (*BigQueryServer, error) {
	s := &BigQueryServer{
		server: server.NewMCPServer(
			"bigquery-server",
			version,
		),
	}

	client, err := NewBigQueryClient(ctx, project, location)
	if err != nil {
		return nil, err
	}
	s.client = client

	s.server.AddTool(mcp.NewTool(string(LIST_TABLES),
		mcp.WithDescription("Get a detailed listing of all tables in a specified dataset."),
		mcp.WithString("dataset",
			mcp.Description("The dataset to list tables from"),
			mcp.Required(),
		),
	), s.handleListTables)

	s.server.AddTool(mcp.NewTool(string(GET_TABLE_SCHEMA),
		mcp.WithDescription("Get the schema of a specified table in a specified dataset."),
		mcp.WithString("dataset",
			mcp.Description("The dataset to get the table schema from"),
			mcp.Required(),
		),
		mcp.WithString("table",
			mcp.Description("The table to get the schema from"),
			mcp.Required(),
		),
	), s.handleGetTableSchema)

	return s, nil
}

type BigQueryServer struct {
	server *server.MCPServer
	client *BigQueryClient
}

func (s *BigQueryServer) Serve() error {
	return server.ServeStdio(s.server)
}

func (s *BigQueryServer) handleListTables(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	dataset, ok := arguments["dataset"].(string)
	if !ok {
		return nil, fmt.Errorf("dataset must be a string")
	}

	tables, err := s.client.ListTables(ctx, dataset)
	if err != nil {
		return nil, err
	}

	var tablesStr string
	for _, table := range tables {
		tablesStr += fmt.Sprintf("- %s\n", table)
	}

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Tables in dataset `%s`:\n\n%s", dataset, tablesStr),
			},
		},
	}, nil
}

func (s *BigQueryServer) handleGetTableSchema(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	dataset, ok := arguments["dataset"].(string)
	if !ok {
		return nil, fmt.Errorf("dataset must be a string")
	}
	table, ok := arguments["table"].(string)
	if !ok {
		return nil, fmt.Errorf("table must be a string")
	}

	schema, err := s.client.GetTableSchema(ctx, dataset, table)
	if err != nil {
		return nil, err
	}

	var schemaStr string
	for _, field := range schema {
		schemaStr += fmt.Sprintf("- %s (%s)\n", field.Name, field.Type)
		if field.Description != "" {
			schemaStr += fmt.Sprintf("  Description: %s\n", field.Description)
		}
		if field.Repeated {
			schemaStr += "  Repeated: true\n"
		}
		if field.Required {
			schemaStr += "  Required: true\n"
		}
		schemaStr += "\n"
	}

	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Schema for table %s in dataset %s:\n\n%s", table, dataset, schemaStr),
			},
		},
	}, nil
}

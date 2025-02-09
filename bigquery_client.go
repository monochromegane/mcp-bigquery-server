package mcp_bigquery_server

import (
	"context"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQueryClient struct {
	Project  string
	Location string
	client   *bigquery.Client
}

func NewBigQueryClient(ctx context.Context, project, location string) (*BigQueryClient, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	client.Location = location
	return &BigQueryClient{
		Project:  project,
		Location: location,
		client:   client,
	}, nil
}

func (c *BigQueryClient) ListTables(ctx context.Context, dataset string) ([]string, error) {
	it := c.client.Dataset(dataset).Tables(ctx)
	tables := []string{}
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		tables = append(tables, t.TableID)
	}
	return tables, nil
}

func (c *BigQueryClient) GetTableSchema(ctx context.Context, dataset, table string) ([]*bigquery.FieldSchema, error) {
	md, err := c.client.Dataset(dataset).Table(table).Metadata(ctx)
	if err != nil {
		return nil, err
	}
	return md.Schema, nil
}

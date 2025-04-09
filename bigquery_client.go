package mcp_bigquery_server

import (
	"context"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQueryClient struct {
	Project string
	client  *bigquery.Client
}

func NewBigQueryClient(ctx context.Context, project string) (*BigQueryClient, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	return &BigQueryClient{
		Project: project,
		client:  client,
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

func (c *BigQueryClient) DryRunQuery(ctx context.Context, query string, dataset string) (*bigquery.JobStatus, error) {
	q := c.client.Query(query)
	q.DefaultProjectID = c.Project
	q.DefaultDatasetID = dataset
	q.DryRun = true
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	return job.LastStatus(), nil
}

package gotest2bq

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
)

type TestEvent struct {
	Branch  string    `bigquery:"branch"`
	Env     string    `bigquery:"env"`
	Commit  string    `bigquery:"commit"`
	GroupId string    `bigquery:"group_id"`
	Time    time.Time `bigquery:"time"`
	Action  string    `bigquery:"action"`
	Package string    `bigquery:"package"`
	Test    string    `bigquery:"test"`
	Elapsed float64   `bigquery:"elapsed"`
	Output  string    `bigquery:"-"`
}

func loadTestLog(args GoTest2BQArgs) ([]*TestEvent, error) {
	testEvents := make([]*TestEvent, 0)

	file, err := os.Open(args.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		testEvent := &TestEvent{}
		err := json.Unmarshal([]byte(line), testEvent)
		if err != nil {
			return nil, err
		}
		testEvent.Branch = args.Branch
		testEvent.Env = args.Env
		testEvent.Commit = args.Commit
		testEvent.GroupId = args.groupId
		testEvents = append(testEvents, testEvent)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return testEvents, nil
}

type GoTest2BQArgs struct {
	Branch   string
	Env      string
	Commit   string
	Filename string
	Project  string
	Dataset  string
	Table    string
	groupId  string
}

func GoTest2BQ(args GoTest2BQArgs) error {
	ctx := context.Background()

	args.groupId = uuid.NewString()

	testEvents, err := loadTestLog(args)
	if err != nil {
		return fmt.Errorf("load test log: %w", err)
	}
	client, err := bigquery.NewClient(ctx, args.Project)
	if err != nil {
		return fmt.Errorf("bigquery client: %w", err)
	}
	defer client.Close()
	schema, err := bigquery.InferSchema(TestEvent{})
	if err != nil {
		return fmt.Errorf("infer schema: %w", err)
	}
	schema = schema.Relax()
	table := client.Dataset(args.Dataset).Table(args.Table)

	tm, err := table.Metadata(ctx)
	if err != nil {
		err := table.Create(ctx, &bigquery.TableMetadata{
			Schema: schema,
			TimePartitioning: &bigquery.TimePartitioning{
				Field: "time",
			},
		},
		)
		if err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	} else {
		_, err = table.Update(ctx, bigquery.TableMetadataToUpdate{
			Schema: schema,
		}, tm.ETag)
	}
	if err != nil {
		return fmt.Errorf("update table: %w", err)
	}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, testEvents); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

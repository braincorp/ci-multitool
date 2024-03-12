package jira

import (
	"fmt"
	"io"
	"os"
	"strings"

	gojira "github.com/andygrunwald/go-jira"
)

type CreateIssueArgs struct {
	Common       *CommonArgs
	Summary      string
	Description  string
	Assignee     string
	Type         string
	Labels       []string
	CustomFields map[string]string
}

func CreateIssue(args *CreateIssueArgs) (string, error) {
	client, err := GetClient(args.Common)
	if err != nil {
		return "", fmt.Errorf("get client: %w", err)
	}

	customFields := make(map[string]interface{})
	for k, v := range args.CustomFields {
		customFields[k] = v
	}

	if args.Type == "Auto" {
		if strings.Contains(strings.ToLower(args.Summary), "fix") {
			args.Type = "Bug"
		} else {
			args.Type = "Task"
		}
	}

	issue := &gojira.Issue{
		Fields: &gojira.IssueFields{
			Project: gojira.Project{
				Key: args.Common.Project,
			},
			Type: gojira.IssueType{
				Name: args.Type,
			},
			Description: args.Description,
			Summary:     args.Summary,
			Labels:      args.Labels,
			Unknowns:    customFields,
		},
	}

	if args.Assignee != "" {
		users, _, err := client.User.Find(args.Assignee)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error finding user: %v\n", err)
		} else {
			issue.Fields.Assignee = &users[0]
		}
	}

	issue, resp, err := client.Issue.Create(issue)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create issue: %w: %s", err, string(body))
	}
	return issue.Key, nil
}

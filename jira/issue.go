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
	Components   []string
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

	if len(args.Components) > 0 {
		components := make([]*gojira.Component, 0, len(args.Components))
		for _, name := range args.Components {
			if name == "" {
				continue
			}
			components = append(components, &gojira.Component{Name: name})
		}
		if len(components) > 0 {
			issue.Fields.Components = components
		}
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

func GetActiveSprint(commonArgs *CommonArgs) (*gojira.Sprint, error) {
	client, err := GetClient(commonArgs)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	sprints, _, err := client.Board.GetAllSprintsWithOptions(commonArgs.Board, &gojira.GetAllSprintsOptions{
		State: "active",
	})
	if err != nil {
		return nil, err
	}
	for _, sprint := range sprints.Values {
		if sprint.State == "active" {
			return &sprint, nil
		}
	}
	return nil, fmt.Errorf("no active sprint found")
}

func AddIssueToSprint(commonArgs *CommonArgs, issueID string, sprintID int) error {
	client, err := GetClient(commonArgs)
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	_, err = client.Sprint.MoveIssuesToSprint(sprintID, []string{issueID})
	return err
}

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/alexgartner-bc/ci-multitool/jira"
	"github.com/spf13/cobra"
)

func init() {
	jiraCmd.AddCommand(jiraCreateIssueCmd)
	jiraCmdF := jiraCmd.PersistentFlags()
	jiraCmdF.String("instance-url", os.Getenv("JIRA_INSTANCE_URL"), "instance url")
	jiraCmdF.String("project", os.Getenv("JIRA_PROJECT"), "project")
	jiraCmdF.String("user", os.Getenv("JIRA_USER"), "user")
	jiraCmdF.String("password", os.Getenv("JIRA_PASSWORD"), "password/token")

	jiraCreateIssueCmdF := jiraCreateIssueCmd.Flags()
	jiraCreateIssueCmdF.StringP("summary", "s", "", "issue summary/title")
	jiraCreateIssueCmdF.StringP("description", "d", "", "issue description")
	jiraCreateIssueCmdF.StringP("assignee", "a", "", "issue assignee")
	jiraCreateIssueCmdF.StringP("type", "t", "", "issue type")
	jiraCreateIssueCmdF.StringSliceP("labels", "l", []string{}, "issue labels")
	jiraCreateIssueCmdF.StringToString("custom", map[string]string{}, "issue custom fields")
	jiraCreateIssueCmdF.StringSlice("components", []string{}, "issue components (repeatable)")
}

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "tools to work with jira",
}

func getCommonArgs(cmd *cobra.Command) (*jira.CommonArgs, error) {
	instanceUrl, _ := cmd.Flags().GetString("instance-url")
	if instanceUrl == "" {
		return nil, errors.New("instance url is required")
	}
	project, _ := cmd.Flags().GetString("project")
	if project == "" {
		return nil, errors.New("project is required")
	}
	user, _ := cmd.Flags().GetString("user")
	if user == "" {
		return nil, errors.New("user is required")
	}
	password, _ := cmd.Flags().GetString("password")
	if password == "" {
		return nil, errors.New("password is required")
	}
	return &jira.CommonArgs{
		InstanceUrl: instanceUrl,
		Project:     project,
		User:        user,
		Password:    password,
	}, nil
}

var jiraCreateIssueCmd = &cobra.Command{
	Use:          "create-issue",
	Short:        "create jira issue",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		commonArgs, err := getCommonArgs(cmd)
		if err != nil {
			return err
		}
		summary, _ := cmd.Flags().GetString("summary")
		if summary == "" {
			return errors.New("summary is required")
		}
		description, _ := cmd.Flags().GetString("description")
		assignee, _ := cmd.Flags().GetString("assignee")
		issueType, _ := cmd.Flags().GetString("type")
		if issueType == "" {
			return errors.New("type is required")
		}
		labels, _ := cmd.Flags().GetStringSlice("labels")
		customFields, _ := cmd.Flags().GetStringToString("custom")
		components, _ := cmd.Flags().GetStringSlice("components")

		key, err := jira.CreateIssue(&jira.CreateIssueArgs{
			Common:       commonArgs,
			Summary:      summary,
			Description:  description,
			Assignee:     assignee,
			Type:         issueType,
			Labels:       labels,
			CustomFields: customFields,
			Components:   components,
		})
		if err != nil {
			return err
		}
		fmt.Println(key)
		return nil
	},
}

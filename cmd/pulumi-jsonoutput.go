package cmd

import (
	"fmt"

	"github.com/alexgartner-bc/multitool/github"
	"github.com/alexgartner-bc/multitool/pulumi/jsonoutput"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var pulumiJSONOutputFlags = struct {
	destinations []string
}{
	destinations: []string{},
}

func init() {
	pulumiJSONOutput.Flags().StringSliceVarP(
		&pulumiJSONOutputFlags.destinations,
		"destinations", "d",
		[]string{},
		"comma separated list of destinations (stdout,gh-comment,gh-pr-trailer)",
	)
	setGithubDefaultArgs(pulumiJSONOutput.Flags())
}

var pulumiJSONOutput = &cobra.Command{
	Use:   "jsonoutput <file>",
	Short: "process the json output from pulumi",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		destinations := pulumiJSONOutputFlags.destinations

		m, err := jsonoutput.NewManagerFromFile(args[0])
		if err != nil {
			return fmt.Errorf("unable to make jsonoutput manager: %w", err)
		}
		summary := m.ShortSummaryString()
		errMessage := m.Error()
		tree := m.TreeString()

		if slices.Contains(destinations, "stdout") {
			fmt.Println("Summary: " + summary + "\n")
			if errMessage != "" {
				fmt.Println(errMessage)
			}
			fmt.Println(tree)
		}
		if slices.Contains(destinations, "gh-pr-trailer") {
			ghSummary := fmt.Sprintf("pulumi output (%s)", summary)
			ghDetails := fmt.Sprintf("```\n%s```", tree)
			if errMessage != "" {
				ghDetails = fmt.Sprintf("```\n%s\n```\n%s", errMessage, ghDetails)
			}
			err = github.SetPRTrailerDetails(ctx,
				githubDefaultArgs.repo,
				githubDefaultArgs.pr,
				ghSummary,
				ghDetails,
				githubDefaultArgs.key,
			)
			if err != nil {
				return fmt.Errorf("unable to set github pr trailer: %w", err)
			}
		}
		return nil
	},
}

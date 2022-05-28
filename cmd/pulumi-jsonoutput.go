package cmd

import (
	"fmt"

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
}

var pulumiJSONOutput = &cobra.Command{
	Use:   "jsonoutput <file>",
	Short: "process the json output from pulumi",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		destinations := pulumiJSONOutputFlags.destinations

		m, err := jsonoutput.NewManagerFromFile(args[0])
		if err != nil {
			return fmt.Errorf("unable to make jsonoutput manager: %w", err)
		}
		summary := m.ShortSummaryString()
		tree := m.TreeString()

		if slices.Contains(destinations, "stdout") {
			fmt.Println("Summary: " + summary + "\n")
			fmt.Println(tree)
		}
		return nil
	},
}

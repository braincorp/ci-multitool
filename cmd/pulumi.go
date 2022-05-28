package cmd

import "github.com/spf13/cobra"

func init() {
	pulumiCmd.AddCommand(pulumiJSONOutput)
}

var pulumiCmd = &cobra.Command{
	Use:   "pulumi",
	Short: "tools to work with pulumi",
}

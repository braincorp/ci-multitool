package cmd

import (
	"errors"

	"github.com/alexgartner-bc/ci-multitool/gotest2bq"
	"github.com/spf13/cobra"
)

func init() {
	fs := gotest2bqCmd.Flags()
	fs.String("project", "", "bigquery project")
	fs.String("dataset", "", "bigquery dataset")
	fs.String("table", "", "bigquery table")
}

var gotest2bqCmd = &cobra.Command{
	Use:   "gotest2bq",
	Short: "ingest go test output into bigquery",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		dataset, _ := cmd.Flags().GetString("dataset")
		table, _ := cmd.Flags().GetString("table")

		if project == "" || dataset == "" || table == "" {
			return errors.New("project, dataset, and table are required")
		}

		filename := args[0]
		err := gotest2bq.GoTest2BQ(filename, project, dataset, table)
		if err != nil {
			return err
		}
		return nil
	},
}

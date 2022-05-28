package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alexgartner-bc/multitool/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var githubDefaultArgs = struct {
	repo string
	pr   int
	key  string
}{}

var githubPRTrailerArgs = struct {
	summary string
}{}

func setGithubDefaultArgs(fs *pflag.FlagSet) {
	fs.StringVar(
		&githubDefaultArgs.repo,
		"repo", "",
		"name of the repo (alexgartner-bc/my-repo)",
	)
	fs.IntVar(
		&githubDefaultArgs.pr,
		"pr", 0,
		"number of the pr (1234)",
	)
	fs.StringVar(
		&githubDefaultArgs.key,
		"key", "default-key",
		"hidden text to embed in the comment to allow updating it later",
	)
}

func init() {
	githubCmd.AddCommand(githubCommentCmd)
	setGithubDefaultArgs(githubCommentCmd.Flags())

	githubCmd.AddCommand(githubPrTrailerCmd)
	setGithubDefaultArgs(githubPrTrailerCmd.Flags())
	githubPrTrailerCmd.Flags().StringVar(
		&githubPRTrailerArgs.summary,
		"summary",
		"",
		"<summary> for the <details>",
	)
}

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "tools to work with github",
}

var githubCommentCmd = &cobra.Command{
	Use:   "comment <file>",
	Short: "comment on github from file (can be - for stdin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var input io.ReadCloser
		var err error
		if args[0] == "-" {
			input = os.Stdin
		} else {
			input, err = os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
		}

		body, err := io.ReadAll(input)
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}

		repo := githubDefaultArgs.repo
		if repo == "" {
			return errors.New("repo must be set")
		}
		prNumber := githubDefaultArgs.pr
		if prNumber == 0 {
			return errors.New("pr must be set")
		}
		key := githubDefaultArgs.key

		err = github.CommentOnIssue(ctx, repo, prNumber, string(body), key)
		return err
	},
}

var githubPrTrailerCmd = &cobra.Command{
	Use:   "pr-trailer <file>",
	Short: "add text to the bottom of the PR from a file (can be - for stdin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var input io.ReadCloser
		var err error
		if args[0] == "-" {
			input = os.Stdin
		} else {
			input, err = os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot open file: %w", err)
			}
		}

		body, err := io.ReadAll(input)
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}

		repo := githubDefaultArgs.repo
		if repo == "" {
			return errors.New("repo must be set")
		}
		prNumber := githubDefaultArgs.pr
		if prNumber == 0 {
			return errors.New("pr must be set")
		}
		key := githubDefaultArgs.key

		err = github.SetPRTrailerDetails(ctx, repo, prNumber, githubPRTrailerArgs.summary, string(body), key)
		return err
	},
}

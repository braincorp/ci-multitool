package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
)

// CommentOnIssue posts a comment on an issue or PR
//
// github.repository => repo (alexgartner-bc/my-repo)
// github.event.issue.number => number
func CommentOnIssue(ctx context.Context, repo string, number int, text string, stickyKey string) error {
	client := getDefaultClient()

	repoParts := strings.Split(repo, "/")

	stickyKeyText := fmt.Sprintf("\n<!-- key %s -->\n", stickyKey)

	text += stickyKeyText

	commentReq := &github.IssueComment{
		Body: &text,
	}

	commentListOptions := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	// try to find existing comment
	existingComments, _, err := client.Issues.ListComments(ctx, repoParts[0], repoParts[1], number, commentListOptions)
	if err != nil {
		return fmt.Errorf("unable to list comments: %w", err)
	}
	for _, comment := range existingComments {
		if strings.Contains(*comment.Body, stickyKeyText) {
			// update existing comment
			_, _, err = client.Issues.EditComment(ctx, repoParts[0], repoParts[1], *comment.ID, commentReq)
			if err != nil {
				return fmt.Errorf("unable to edit comment: %w", err)
			}
			return nil
		}
	}

	_, _, err = client.Issues.CreateComment(ctx, repoParts[0], repoParts[1], number, commentReq)
	if err != nil {
		return fmt.Errorf("unable to create comment: %w", err)
	}
	return nil
}

func CommentOnCommit(ctx context.Context, repo string, sha string, text string, stickyKey string) error {
	client := getDefaultClient()

	repoParts := strings.Split(repo, "/")

	stickyKeyText := fmt.Sprintf("\n<!-- key %s -->\n", stickyKey)

	text += stickyKeyText

	commentReq := &github.RepositoryComment{
		Body: &text,
	}

	// try to find existing comment
	existingComments, _, err := client.Repositories.ListCommitComments(ctx, repoParts[0], repoParts[1], sha, nil)
	if err != nil {
		return fmt.Errorf("unable to list comments: %w", err)
	}
	for _, comment := range existingComments {
		if strings.Contains(*comment.Body, stickyKeyText) {
			// update existing comment
			_, _, err = client.Repositories.UpdateComment(ctx, repoParts[0], repoParts[1], *comment.ID, commentReq)
			if err != nil {
				return fmt.Errorf("unable to edit comment: %w", err)
			}
			return nil
		}
	}

	_, _, err = client.Repositories.CreateComment(ctx, repoParts[0], repoParts[1], sha, commentReq)
	if err != nil {
		return fmt.Errorf("unable to create comment: %w", err)
	}
	return nil
}

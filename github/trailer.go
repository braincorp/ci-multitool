package github

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// SetPRTrailerDetails update the a PR and sets some text at the bottom. Might be better than a comment because it doesn't cause a notification.
// If summary is set, use <details>. Else use <span>
//
// github.repository => repo (alexgartner-bc/my-repo)
// github.event.issue.number => number
func SetPRTrailerDetails(ctx context.Context, repo string, number int, summary string, details string, stickyKey string) error {
	client := getDefaultClient()

	repoParts := strings.Split(repo, "/")

	tag := "span"
	if summary != "" {
		tag = "details"
	}

	openingTag := fmt.Sprintf("<%s id=\"%s\">", tag, stickyKey)
	closingTag := fmt.Sprintf("</%s>", tag)

	summaryTag := ""
	if summary != "" {
		summaryTag = fmt.Sprintf("<summary>%s</summary>", summary)
	}

	text := fmt.Sprintf("\n%s%s\n\n%s\n\n%s", openingTag, summaryTag, details, closingTag)

	pr, _, err := client.PullRequests.Get(ctx, repoParts[0], repoParts[1], number)
	if err != nil {
		return fmt.Errorf("unable to get PR: %w", err)
	}

	body := ""
	if pr.Body != nil {
		body = *pr.Body
	}

	if strings.Contains(body, openingTag) {
		re := regexp.MustCompile(fmt.Sprintf("(?ms)\n%s.+?%s", openingTag, closingTag))
		body = re.ReplaceAllString(body, text)
	} else {
		body += text
	}

	pr.Body = &body

	_, _, err = client.PullRequests.Edit(ctx, repoParts[0], repoParts[1], number, pr)
	if err != nil {
		return fmt.Errorf("unable to edit PR: %w", err)
	}
	return nil
}

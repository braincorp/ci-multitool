package github

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// SetPRTrailerDetails update the a PR and sets some text at the bottom. Might be better than a comment because it doesn't cause a notification.
//
// github.repository => repo (alexgartner-bc/my-repo)
// github.event.issue.number => number
func SetPRTrailerDetails(ctx context.Context, repo string, number int, summary string, details string, stickyKey string) error {
	client := getDefaultClient()

	repoParts := strings.Split(repo, "/")

	detailsOpeningTag := fmt.Sprintf("<details key=\"%s\">", stickyKey)

	text := fmt.Sprintf("\n%s<summary>%s</summary>\n\n%s\n\n</details>", detailsOpeningTag, summary, details)

	pr, _, err := client.PullRequests.Get(ctx, repoParts[0], repoParts[1], number)
	if err != nil {
		return fmt.Errorf("unable to get PR: %w", err)
	}

	body := ""
	if pr.Body != nil {
		body = *pr.Body
	}

	if strings.Contains(body, detailsOpeningTag) {
		re := regexp.MustCompile(fmt.Sprintf("(?ms)%s.+?</details>", detailsOpeningTag))
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

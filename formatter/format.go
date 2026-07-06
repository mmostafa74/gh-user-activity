package formatter

import (
	"fmt"
	"strings"

	types "gh-user-activity/client"
)

func FormatEvent(e types.Event) string {
	repoName := e.Repo.Name
	switch e.Type {
	case "PushEvent":
		branch := shortBranch(e.Payload.Ref)
		if branch != "" {
			return fmt.Sprintf("Pushed to %s in %s", branch, repoName)
		}
		return fmt.Sprintf("Pushed to %s", repoName)

	case "IssuesEvent":
		issue := e.Payload.Issue
		if issue != nil {
			return fmt.Sprintf("%s issue #%d in %s", capitalize(e.Payload.Action), issue.Number, repoName)
		}
		return fmt.Sprintf("%s an issue in %s", e.Payload.Action, repoName)

	case "WatchEvent":
		return fmt.Sprintf("Starred %s", repoName)

	case "ForkEvent":
		return fmt.Sprintf("Forked %s", repoName)

	case "CreateEvent":
		refType := e.Payload.RefType
		ref := e.Payload.Ref
		if ref != "" {
			return fmt.Sprintf("Created %s %s in %s", refType, ref, repoName)
		}
		return fmt.Sprintf("Created %s in %s", refType, repoName)

	case "DeleteEvent":
		return fmt.Sprintf("Deleted %s %s in %s", e.Payload.RefType, e.Payload.Ref, repoName)

	case "PullRequestEvent":
		pr := e.Payload.PullRequest
		if pr != nil {
			return fmt.Sprintf("%s PR #%d in %s", capitalize(e.Payload.Action), pr.Number, repoName)
		}
		return fmt.Sprintf("%s a PR in %s", e.Payload.Action, repoName)

	case "IssueCommentEvent":
		issue := e.Payload.Issue
		if issue != nil {
			return fmt.Sprintf("Commented on issue #%d in %s", issue.Number, repoName)
		}
		return fmt.Sprintf("Commented on an issue in %s", repoName)

	case "ReleaseEvent":
		return fmt.Sprintf("Published release in %s", repoName)

	case "MemberEvent":
		member := e.Payload.Member
		if member != nil {
			return fmt.Sprintf("%s %s to %s", capitalize(e.Payload.Action), member.Login, repoName)
		}
		return fmt.Sprintf("%s a member in %s", e.Payload.Action, repoName)

	case "PublicEvent":
		return fmt.Sprintf("Made %s public", repoName)

	case "GollumEvent":
		return fmt.Sprintf("Updated wiki in %s", repoName)

	case "PullRequestReviewEvent":
		pr := e.Payload.PullRequest
		if pr != nil {
			return fmt.Sprintf("%s review on PR #%d in %s", capitalize(e.Payload.Action), pr.Number, repoName)
		}
		return fmt.Sprintf("%s a review on a PR in %s", e.Payload.Action, repoName)

	case "PullRequestReviewCommentEvent":
		pr := e.Payload.PullRequest
		if pr != nil {
			return fmt.Sprintf("Commented on PR #%d in %s", pr.Number, repoName)
		}
		return fmt.Sprintf("Commented on a PR in %s", repoName)

	default:
		return fmt.Sprintf("%s in %s", e.Type, repoName)
	}
}

func shortBranch(ref string) string {
	if strings.HasPrefix(ref, "refs/heads/") {
		return ref[len("refs/heads/"):]
	}
	if strings.HasPrefix(ref, "refs/tags/") {
		return ref[len("refs/tags/"):]
	}
	return ref
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return string(s[0]-32) + s[1:]
}

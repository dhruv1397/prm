package client

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/google/go-github/v64/github"
	"net/url"
	"strings"
)

type GithubClient struct {
	user   string
	client *github.Client
}

var _ Client = (*GithubClient)(nil)

func NewGithubClient(ctx context.Context, client *github.Client) (*GithubClient, error) {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error fetching authenticated user: %v", err)
	}

	return &GithubClient{
		user:   user.GetLogin(),
		client: client,
	}, nil
}

func (g *GithubClient) GetOpenPullRequests(ctx context.Context) ([]*types.PullRequest, error) {
	var openPRs = make([]*types.PullRequest, 0)

	query := fmt.Sprintf("state:open author:%s type:pr", g.user)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	result, _, err := g.client.Search.Issues(ctx, query, opts)
	if err != nil {
		return openPRs, fmt.Errorf("error executing search query: %v", err)
	}

	for _, issue := range result.Issues {
		openPR, err2 := g.getPRDetails(ctx, issue)
		if err2 != nil {
			return openPRs, err2
		}

		openPRs = append(openPRs, openPR)

	}
	return openPRs, nil
}

func (g *GithubClient) getPRDetails(
	ctx context.Context,
	issue *github.Issue,
) (*types.PullRequest, error) {
	owner, repo, parseErr := parseGithubURL(*issue.HTMLURL)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing query: %v", parseErr)
	}

	pr, _, err := g.client.PullRequests.Get(ctx, owner, repo, *issue.Number)
	if err != nil {
		return nil, fmt.Errorf("error fetching PR details: %v", err)
	}

	reviews, _, err := g.client.PullRequests.ListReviews(ctx, owner, repo, *issue.Number, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching reviews: %v", err)
	}

	approvedMap := map[string]bool{}
	commentedMap := map[string]bool{}
	changesRequestedMap := map[string]bool{}

	for _, review := range reviews {
		userName := *(review.User.Login)
		if review.GetState() == "APPROVED" {
			approvedMap[userName] = true
		}
		if review.GetState() == "COMMENTED" {
			commentedMap[userName] = true
		}
		if review.GetState() == "CHANGES_REQUESTED" {
			changesRequestedMap[userName] = true
		}
	}

	approved := make([]string, 0, len(approvedMap))
	for k := range approvedMap {
		approved = append(approved, k)
	}

	commented := make([]string, 0, len(commentedMap))
	for k := range commentedMap {
		commented = append(commented, k)
	}

	changesRequested := make([]string, 0, len(changesRequestedMap))
	for k := range changesRequestedMap {
		changesRequested = append(changesRequested, k)
	}

	return &types.PullRequest{
		Title:            *pr.Title,
		Number:           *pr.Number,
		SCMProviderType:  "github",
		URL:              issue.GetHTMLURL(),
		Mergeable:        pr.Mergeable != nil && *pr.Mergeable,
		Approved:         approved,
		Commented:        commented,
		ChangesRequested: changesRequested,
	}, nil

}

func parseGithubURL(githubURL string) (string, string, error) {
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid GitHub URL format")
	}

	owner := parts[1]
	repo := parts[2]

	return owner, repo, nil
}

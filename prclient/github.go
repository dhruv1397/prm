package prclient

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/dhruv1397/pr-monitor/util"
	"github.com/google/go-github/v64/github"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type GithubPRClient struct {
	user   *types.User
	client *github.Client
}

var _ PRClient = (*GithubPRClient)(nil)

func NewGithubPRClient(user *types.User, client *github.Client) (*GithubPRClient, error) {
	return &GithubPRClient{
		user:   user,
		client: client,
	}, nil
}

func (g *GithubPRClient) GetPullRequests(
	ctx context.Context,
	state string,
	transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
) ([]*types.PrintablePullRequest, error) {
	var openPrintablePRs = make([]*types.PrintablePullRequest, 0)
	var githubState = ""
	if state == "closed" {
		githubState = "state:closed is:unmerged"
	} else if state == "merged" {
		githubState = "state:closed is:merged"
	} else if state == "open" {
		githubState = "state:open"
	} else if state == "all" {
		githubState = ""
	}

	query := fmt.Sprintf("%s author:%s type:pr", githubState, g.user.Name)

	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 500},
	}

	result, _, err := g.client.Search.Issues(ctx, query, opts)
	if err != nil {
		return openPrintablePRs, fmt.Errorf("error fetching github PRs for user %s: %w", g.user.Name, err)
	}

	var prMutex sync.Mutex
	var errMutex sync.Mutex
	var wg sync.WaitGroup

	errCh := make(chan error, 500)
	respCh := make(chan *types.PrintablePullRequest, 500)

	for _, issue := range result.Issues {
		wg.Add(1)

		go func(issue *github.Issue) {
			defer wg.Done()

			openPrintablePR, err := g.getPRDetails(ctx, issue, transformationFn)
			if err != nil {
				errCh <- err
				return
			}

			respCh <- openPrintablePR
		}(issue)

	}

	go func() {
		wg.Wait()
		close(respCh)
		close(errCh)
	}()

	var errs []error

	for respCh != nil || errCh != nil {
		select {
		case resp, ok := <-respCh:
			if !ok {
				respCh = nil
			} else {
				prMutex.Lock()
				openPrintablePRs = append(openPrintablePRs, resp)
				prMutex.Unlock()
			}
		case errValue, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				errMutex.Lock()
				errs = append(errs, errValue)
				errMutex.Unlock()
			}
		}
	}

	if len(errs) > 0 {
		return openPrintablePRs, fmt.Errorf("errors encountered:\n%v", util.FormatErrors(errs))
	}

	return openPrintablePRs, nil
}

func (g *GithubPRClient) getPRDetails(
	ctx context.Context,
	issue *github.Issue,
	transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
) (*types.PrintablePullRequest, error) {
	owner, repo, parseErr := parseGithubURL(*issue.HTMLURL)
	if parseErr != nil {
		return nil, fmt.Errorf("error parsing github PR URL %s: %w", *issue.HTMLURL, parseErr)
	}

	pr, _, err := g.client.PullRequests.Get(ctx, owner, repo, *issue.Number)
	if err != nil {
		return nil, fmt.Errorf("error fetching PR details for %s: %w", *issue.HTMLURL, err)
	}

	reviews, _, err := g.client.PullRequests.ListReviews(ctx, owner, repo, *issue.Number, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching PR reviews for %s: %w", *issue.HTMLURL, err)
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

	state := *pr.State
	var mergeable = strconv.FormatBool(pr.Mergeable != nil && *pr.Mergeable)
	if *pr.Merged {
		state = "merged"
		mergeable = "-"
	}

	openPR := &types.PullRequest{
		Title:            *pr.Title,
		Number:           *pr.Number,
		SCMProviderType:  "github",
		URL:              issue.GetHTMLURL(),
		State:            state,
		Mergeable:        mergeable,
		Approved:         approved,
		Commented:        commented,
		RequestedChanges: changesRequested,
	}

	return transformationFn(openPR), nil
}

func parseGithubURL(githubURL string) (string, string, error) {
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return "", "", fmt.Errorf("error parsing github URL %s: %w", githubURL, err)
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid GitHub URL %s", parsedURL)
	}

	owner := parts[1]
	repo := parts[2]

	return owner, repo, nil
}

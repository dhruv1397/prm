package prclient

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/harness"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/dhruv1397/pr-monitor/util"
	"net/http"
	"strconv"
	"sync"
)

type HarnessPRClient struct {
	httpClient *http.Client
	host       string
	user       *types.User
	repos      []*types.Repo
}

var _ PRClient = (*HarnessPRClient)(nil)

func NewHarnessPRClient(host string, user *types.User, repos []*types.Repo) (*HarnessPRClient, error) {
	return &HarnessPRClient{
		httpClient: http.DefaultClient,
		host:       host,
		user:       user,
		repos:      repos,
	}, nil
}

func (h *HarnessPRClient) GetPullRequests(
	ctx context.Context,
	state string,
	transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
) ([]*types.PrintablePullRequest, error) {
	var allPullRequests []*types.PrintablePullRequest
	var prMutex sync.Mutex
	var errMutex sync.Mutex

	var wg sync.WaitGroup
	errChan := make(chan error, len(h.repos)*500)
	prChan := make(chan *types.PrintablePullRequest, len(h.repos)*500)

	for _, repo := range h.repos {
		wg.Add(1)
		go func(repo *types.Repo) {
			defer wg.Done()

			prs, err := h.getPRs(ctx, repo, state)
			if err != nil {
				errChan <- err
				return
			}

			var prWg sync.WaitGroup
			for _, pr := range prs {
				prWg.Add(1)
				go func(pr *types.PRData) {
					defer prWg.Done()

					prActivities, err := h.getPRActivities(ctx, repo, pr)
					if err != nil {
						errChan <- err
						return
					}

					approvedMap := map[string]bool{}
					commentedMap := map[string]bool{}
					changesRequestedMap := map[string]bool{}
					for _, prActivity := range prActivities {
						if (prActivity.Type == "code-comment" || prActivity.Type == "comment") &&
							!commentedMap[prActivity.PRActivityAuthor.DisplayName] {
							commentedMap[prActivity.PRActivityAuthor.DisplayName] = true
						} else if prActivity.Type == "review-submit" && *prActivity.PRActivityDecision.Decision == "approved" &&
							!approvedMap[prActivity.PRActivityAuthor.DisplayName] {
							approvedMap[prActivity.PRActivityAuthor.DisplayName] = true
						} else if prActivity.Type == "review-submit" && *prActivity.PRActivityDecision.Decision == "changereq" &&
							!changesRequestedMap[prActivity.PRActivityAuthor.DisplayName] {
							changesRequestedMap[prActivity.PRActivityAuthor.DisplayName] = true
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

					url := h.getHarnessPRURL(pr.Number, repo)

					mergeable := "-"
					if pr.State != "merged" {
						mergeable = strconv.FormatBool(pr.MergeCheckStatus == "mergeable")
					}

					currentPullRequest := &types.PullRequest{
						Number:           pr.Number,
						Title:            pr.Title,
						SCMProviderType:  "harness",
						URL:              url,
						Approved:         approved,
						Commented:        commented,
						RequestedChanges: changesRequested,
						Mergeable:        mergeable,
						State:            pr.State,
					}

					prChan <- transformationFn(currentPullRequest)
				}(pr)
			}

			prWg.Wait()
		}(repo)
	}

	go func() {
		wg.Wait()
		close(prChan)
		close(errChan)
	}()

	var errs []error

	for prChan != nil || errChan != nil {
		select {
		case resp, ok := <-prChan:
			if !ok {
				prChan = nil
			} else {
				prMutex.Lock()
				allPullRequests = append(allPullRequests, resp)
				prMutex.Unlock()
			}
		case errValue, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errMutex.Lock()
				errs = append(errs, errValue)
				errMutex.Unlock()
			}
		}
	}

	if len(errs) > 0 {
		return allPullRequests, fmt.Errorf("errors encountered:\n%v", util.FormatErrors(errs))
	}

	return allPullRequests, nil
}

func (h *HarnessPRClient) getHarnessPRURL(prNumber int, repo *types.Repo) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%d", h.host, "/ng/account/", repo.AccountIdentifier,
		"/module/code/orgs/", repo.OrgIdentifier, "/projects/", repo.ProjectIdentifier, "/repos/",
		repo.RepoIdentifier, "/pulls/", prNumber)
}

func (h *HarnessPRClient) getPRs(ctx context.Context, repo *types.Repo, state string) ([]*types.PRData, error) {
	var prs = make([]*types.PRData, 0)
	apiURL := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%d%s", h.host, "/code/api/v1/repos/", repo.RepoIdentifier,
		"/pullreq?accountIdentifier=", repo.AccountIdentifier, "&orgIdentifier=", repo.OrgIdentifier,
		"&projectIdentifier=", repo.ProjectIdentifier, "&state=", state, "&page=0&limit=500&created_by=",
		h.user.PrincipalID, "&order=asc")
	err := harness.Get(ctx, h.httpClient, h.user.PAT, apiURL, &prs)
	if err != nil {
		return prs, fmt.Errorf("error fetching PRs for repo %s: %w", repo.RepoIdentifier, err)
	}
	return prs, nil
}

func (h *HarnessPRClient) getPRActivities(
	ctx context.Context,
	repo *types.Repo,
	pr *types.PRData,
) ([]*types.PRActivity, error) {
	var prActivities = make([]*types.PRActivity, 0)
	apiURL := fmt.Sprintf("%s%s%s%s%d%s%s%s%s%s%s%s", h.host, "/code/api/v1/repos/", repo.RepoIdentifier,
		"/pullreq/", pr.Number, "/activities?accountIdentifier=", repo.AccountIdentifier, "&orgIdentifier=",
		repo.OrgIdentifier, "&projectIdentifier=", repo.ProjectIdentifier, "&type=code-comment&type=comment&type=review-submit")
	err := harness.Get(ctx, h.httpClient, h.user.PAT, apiURL, &prActivities)
	if err != nil {
		return prActivities, fmt.Errorf("error fetching PR activities for %s: %w",
			h.getHarnessPRURL(pr.Number, repo), err)
	}
	return prActivities, nil
}

package prclient

import (
	"context"
	"github.com/dhruv1397/pr-monitor/types"
	"net/http"
)

type HarnessPRClient struct {
	httpClient *http.Client
	repo       *types.HarnessRepo
	user       *types.HarnessUser
}

var _ PRClient = (*HarnessPRClient)(nil)

func NewHarnessPRClient(user *types.HarnessUser, repo *types.HarnessRepo) (*HarnessPRClient, error) {
	return &HarnessPRClient{
		httpClient: http.DefaultClient,
		user:       user,
		repo:       repo,
	}, nil
}

func (h *HarnessPRClient) GetPullRequests(
	ctx context.Context,
	state *string,
	transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
) ([]*types.PrintablePullRequest, error) {
	//TODO implement me
	panic("implement me")
}

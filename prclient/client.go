package prclient

import (
	"context"
	"github.com/dhruv1397/pr-monitor/types"
)

type PRClient interface {
	// GetOpenPullRequests returns a list of all open PRs for the target SCM Provider.
	GetPullRequests(
		ctx context.Context,
		state *string,
		transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
	) ([]*types.PrintablePullRequest, error)
}

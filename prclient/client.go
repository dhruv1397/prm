package prclient

import (
	"context"
	"github.com/dhruv1397/pr-monitor/types"
)

type PRClient interface {
	GetPullRequests(
		ctx context.Context,
		state string,
		transformationFn func(*types.PullRequest) *types.PrintablePullRequest,
	) ([]*types.PullRequestResponse, error)
}

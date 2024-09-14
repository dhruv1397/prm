package client

import (
	"context"
	"github.com/dhruv1397/pr-monitor/types"
)

type Client interface {
	// GetOpenPullRequests returns a list of all open PRs for the target SCM Provider.
	GetOpenPullRequests(ctx context.Context) ([]*types.PullRequest, error)
}

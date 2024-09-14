package client

type Client interface {
	// GetOpenPullRequests returns a list of all open PRs for the target SCM Provider.
	GetOpenPullRequests (ctx context) ([]*types.PullRequest, error)
}
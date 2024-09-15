package scmclient

import (
	"context"
	"fmt"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/google/go-github/v64/github"
)

type GithubSCMClient struct {
	client *github.Client
}

func NewGithubSCMClient(client *github.Client) (*GithubSCMClient, error) {
	return &GithubSCMClient{
		client: client,
	}, nil
}

func (c *GithubSCMClient) GetUser(ctx context.Context) (*types.User, error) {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error fetching authenticated user: %v", err)
	}
	return &types.User{
		Name: user.GetLogin(),
	}, nil
}

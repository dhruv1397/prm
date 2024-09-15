package clientbuilder

import (
	"context"
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/scmclient"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

func GetGithubSCMClient(ctx context.Context, pat string) (*scmclient.GithubSCMClient, error) {
	return scmclient.NewGithubSCMClient(getGithubClientWithPAT(ctx, pat))
}

func GetGithubPRClient(ctx context.Context, user *types.User) (prclient.PRClient, error) {
	return prclient.NewGithubPRClient(user, getGithubClientWithPAT(ctx, user.PAT))
}

func getGithubClientWithPAT(ctx context.Context, pat string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	newClient := github.NewClient(tc)
	return newClient
}

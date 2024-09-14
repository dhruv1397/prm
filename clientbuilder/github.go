package clientbuilder

import (
	"context"
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/scmclient"
	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

func GetGithubSCMClient(ctx context.Context, pat string) (*scmclient.GithubSCMClient, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	newClient := github.NewClient(tc)
	return scmclient.NewGithubSCMClient(newClient)
}

func GetGithubPRClient(ctx context.Context, pat string) (prclient.PRClient, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	newClient := github.NewClient(tc)
	return prclient.NewGithubPRClient(ctx, newClient)
}

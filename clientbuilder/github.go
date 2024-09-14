package clientbuilder

import (
	"context"
	"github.com/dhruv1397/pr-monitor/client"
	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context, pat string) (client.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	newClient := github.NewClient(tc)
	return client.NewGithubClient(ctx, newClient)
}

package clientbuilder

import (
	"context"
	"fmt"
	github "github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
)

type GithubClientBuilder struct {}

func (g *GithubClientBuilder) GetClient (ctx context, pat string) (client.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	newClient := github.NewClient(tc)

	return client.NewGithubClient(ctx, newClient)
}
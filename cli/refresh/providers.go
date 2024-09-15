package refresh

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
	"time"
)

type providersCommand struct {
	name         string
	providerType string
}

func (c *providersCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	str := store.NewSCMProviderImpl()

	providers, err := str.List(c.name, c.providerType)
	if err != nil {
		return err
	}

	var updatedProviders = make([]types.SCMProvider, 0)

	for _, provider := range providers {
		var currentProvider = *provider
		if currentProvider.Type == "github" {
			scmClient, err := clientbuilder.GetGithubSCMClient(ctx, currentProvider.User.PAT)
			if err != nil {
				return err
			}
			gituhbUser, err := scmClient.GetUser(ctx)
			if err != nil {
				return err
			}
			gituhbUser.PAT = currentProvider.User.PAT
			currentProvider.User = gituhbUser
		} else if c.providerType == "harness" {
			scmClient, err := clientbuilder.GetHarnessSCMClient(currentProvider.Host, currentProvider.User.PAT)
			if err != nil {
				return err
			}
			harnessUser, err := scmClient.GetUser(ctx)
			if err != nil {
				return err
			}
			currentProvider.User = harnessUser
			repos, err := scmClient.GetRepos(ctx)
			if err != nil {
				return err
			}
			currentProvider.Repos = repos
		} else {
			return fmt.Errorf("unknown provider type: %s", c.providerType)
		}
		updatedProviders = append(updatedProviders, currentProvider)
	}

	err = str.UpdateBulk(updatedProviders)
	if err != nil {
		return err
	}
	return nil
}

func registerProviders(app *kingpin.CmdClause) {
	c := &providersCommand{}

	cmd := app.Command("providers", "refresh all the SCM providers").Default().Action(c.run)
	cmd.Flag("name", "name of the SCM provider").StringVar(&c.name)
	cmd.Flag("type", "type of the SCM provider").StringVar(&c.name)

}

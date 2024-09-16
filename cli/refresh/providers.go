package refresh

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
	"github.com/dhruv1397/pr-monitor/util"
	"sync"
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

	var updatedProviders []types.SCMProvider
	var errs []error

	var wg sync.WaitGroup
	respCh := make(chan types.SCMProvider)
	errCh := make(chan error)

	var providersMu sync.Mutex
	var errorsMu sync.Mutex

	for _, provider := range providers {
		wg.Add(1)
		go func(provider *types.SCMProvider) {
			defer wg.Done()

			var currentProvider = *provider
			if currentProvider.Type == "github" {
				scmClient, err := clientbuilder.GetGithubSCMClient(ctx, currentProvider.User.PAT)
				if err != nil {
					errCh <- err
					return
				}
				gituhbUser, err := scmClient.GetUser(ctx)
				if err != nil {
					errCh <- err
					return
				}
				gituhbUser.PAT = currentProvider.User.PAT
				currentProvider.User = gituhbUser

			} else if currentProvider.Type == "harness" {
				scmClient, err := clientbuilder.GetHarnessSCMClient(currentProvider.Host, currentProvider.User.PAT)
				if err != nil {
					errCh <- err
					return
				}
				harnessUser, err := scmClient.GetUser(ctx)
				if err != nil {
					errCh <- err
					return
				}
				currentProvider.User = harnessUser
				repos, err := scmClient.GetRepos(ctx)
				if err != nil {
					errCh <- err
					return
				}
				currentProvider.Repos = repos

			} else {
				errCh <- fmt.Errorf("unknown provider type: %s", currentProvider.Type)
				return
			}

			respCh <- currentProvider
		}(provider)
	}

	go func() {
		wg.Wait()
		close(respCh)
		close(errCh)
	}()

	for respCh != nil || errCh != nil {
		select {
		case resp, ok := <-respCh:
			if !ok {
				respCh = nil
			} else {
				providersMu.Lock()
				updatedProviders = append(updatedProviders, resp)
				providersMu.Unlock()
			}
		case errValue, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				errorsMu.Lock()
				errs = append(errs, errValue)
				errorsMu.Unlock()
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors encountered:\n%v", util.FormatErrors(errs))
	}

	err = str.UpdateBulk(updatedProviders)
	if err != nil {
		return err
	}

	return nil
}

func registerProviders(app *kingpin.CmdClause) {
	c := &providersCommand{}

	cmd := app.Command(cli.SubcommandProviders, cli.SubcommandRefreshProvidersHelpText).Default().Action(c.run)

	cmd.Flag(cli.FlagName, cli.FlagNameHelpText).Short(cli.FlagNameShort).StringVar(&c.name)

	cmd.Flag(cli.FlagType, cli.FlagTypeHelpText).Short(cli.FlagTypeShort).StringVar(&c.providerType)

}

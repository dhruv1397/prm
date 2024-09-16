package add

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
	"net/url"
	"strings"
	"time"
)

type providerCommand struct {
	name         string
	providerType string
	host         string
	pat          string
}

func (c *providerCommand) run(*kingpin.ParseContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	str := store.NewSCMProviderImpl()

	existingProviders, err := str.List(c.providerType, c.name)
	if err != nil {
		return err
	}

	if len(existingProviders) != 0 {
		return fmt.Errorf("SCM provider %s already exists", c.name)
	}

	host, err := url.Parse(c.host)
	if err != nil {
		return err
	}

	if host.Scheme == "" {
		host.Scheme = "https"
	}

	pat, err := promptForSecret("Enter the PAT (Personal Access Token):")
	if err != nil {
		return err
	}

	c.host = strings.TrimSuffix(host.String(), "/")
	newProvider := &types.SCMProvider{
		Type:    c.providerType,
		Name:    c.name,
		Host:    c.host,
		User:    nil,
		Repos:   nil,
		Updated: time.Now().UnixMilli(),
		Created: time.Now().UnixMilli(),
	}

	if c.providerType == "github" {
		scmClient, err := clientbuilder.GetGithubSCMClient(ctx, pat)
		if err != nil {
			return err
		}
		gituhbUser, err := scmClient.GetUser(ctx)
		if err != nil {
			return err
		}
		gituhbUser.PAT = pat
		newProvider.User = gituhbUser
	} else if c.providerType == "harness" {
		scmClient, err := clientbuilder.GetHarnessSCMClient(newProvider.Host, pat)
		if err != nil {
			return err
		}
		harnessUser, err := scmClient.GetUser(ctx)
		if err != nil {
			return err
		}
		newProvider.User = harnessUser
		repos, err := scmClient.GetRepos(ctx)
		if err != nil {
			return err
		}
		newProvider.Repos = repos
	} else {
		return fmt.Errorf("unknown provider type: %s", c.providerType)
	}

	err = str.Create(*newProvider)
	if err != nil {
		return err
	}
	return nil
}

func registerProvider(app *kingpin.CmdClause) {
	c := &providerCommand{}

	cmd := app.Command(cli.SubcommandProvider, cli.SubcommandAddProviderHelpText).Action(c.run)

	cmd.Arg(cli.ArgName, cli.ArgNameHelpText).Required().StringVar(&c.name)

	cmd.Flag(cli.FlagType, cli.FlagTypeHelpText).Required().StringVar(&c.providerType)

	cmd.Flag(cli.FlagHost, cli.FlagHostHelpText).Required().StringVar(&c.host)
}

func promptForSecret(promptText string) (string, error) {
	var result string
	prompt := &survey.Password{
		Message: promptText,
	}
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", fmt.Errorf("failed to prompt for secret: %s", err)
	}
	return result, nil
}

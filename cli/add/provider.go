package add

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/clientbuilder"
	"github.com/dhruv1397/pr-monitor/store"
	"github.com/dhruv1397/pr-monitor/types"
	"golang.org/x/term"
	"net/url"
	"os"
	"strings"
	"syscall"
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

	pat := promptForSecret()

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

	cmd.Flag(cli.FlagType, cli.FlagTypeHelpText).Short(cli.FlagTypeShort).Required().StringVar(&c.providerType)

	cmd.Flag(cli.FlagHost, cli.FlagHostHelpText).Short(cli.FlagHostShort).Required().StringVar(&c.host)
}

func promptForSecret() string {
	fmt.Println("Enter the PAT (Personal Access Token):")
	patBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	return string(patBytes)
}

package list

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/store"
)

type providersCommand struct {
	providerType string
	providerName string
}

func (c *providersCommand) run(*kingpin.ParseContext) error {
	str := store.NewSCMProviderImpl()
	providers, err := str.List(c.providerType, c.providerName)
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}
	if len(providers) == 0 {
		fmt.Println("No providers found!")
		return nil
	}
	fmt.Println(fmt.Sprintf("%-4s\t%-10s\t%-10s\t%-20s", "#", "Name", "Type", "Host"))
	for i, provider := range providers {
		fmt.Println(fmt.Sprintf("%-4d\t%-10s\t%-10s\t%-20s", i, provider.Name, provider.Type, provider.Host))
	}
	return nil
}

func registerProviders(app *kingpin.CmdClause) {
	c := &providersCommand{}

	cmd := app.Command(cli.SubcommandProviders, cli.SubcommandListProvidersHelpText).Action(c.run)

	cmd.Flag(cli.FlagType, cli.FlagTypeHelpText).Short(cli.FlagTypeShort).StringVar(&c.providerType)

	cmd.Flag(cli.FlagName, cli.FlagNameHelpText).Short(cli.FlagNameShort).StringVar(&c.providerName)
}

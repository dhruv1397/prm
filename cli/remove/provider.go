package remove

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/prm/cli"
	"github.com/dhruv1397/prm/store"
)

type providerCommand struct {
	name string
}

func (c *providerCommand) run(*kingpin.ParseContext) error {
	str := store.NewSCMProviderImpl()
	err := str.Delete(c.name)
	if err != nil {
		return err
	}
	return nil
}

func registerProvider(app *kingpin.CmdClause) {
	c := &providerCommand{}

	cmd := app.Command(cli.SubcommandProvider, cli.SubcommandRemoveProviderHelpText).Action(c.run)

	cmd.Arg(cli.ArgName, cli.ArgNameHelpText).Required().StringVar(&c.name)
}

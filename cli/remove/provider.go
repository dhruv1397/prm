package remove

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/store"
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

	cmd := app.Command("provider", "remove an SCM provider").Action(c.run)

	cmd.Arg("name", "name of the SCM provider").Required().StringVar(&c.name)
}

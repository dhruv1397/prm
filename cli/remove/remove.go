package remove

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/prm/cli"
)

func Register(app *kingpin.Application) {
	cmd := app.Command(cli.CommandRemove, cli.CommandRemoveHelpText)
	registerProvider(cmd)
}

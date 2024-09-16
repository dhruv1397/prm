package list

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/prm/cli"
)

func Register(app *kingpin.Application) {
	cmd := app.Command(cli.CommandList, cli.CommandListHelpText)
	registerPRs(cmd)
	registerProviders(cmd)
}

package add

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
)

func Register(app *kingpin.Application) {
	cmd := app.Command(cli.CommandAdd, cli.CommandAddHelpText)
	registerProvider(cmd)
}

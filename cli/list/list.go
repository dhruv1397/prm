package list

import "github.com/alecthomas/kingpin/v2"

func Register(app *kingpin.Application) {
	cmd := app.Command("list", "list pull requests and SCM providers")
	registerPRs(cmd)
	registerProviders(cmd)
}

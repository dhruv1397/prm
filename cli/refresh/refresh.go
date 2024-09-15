package refresh

import "github.com/alecthomas/kingpin/v2"

func Register(app *kingpin.Application) {
	cmd := app.Command("refresh", "refresh fetched values ie repos, user-name for all the SCM providers")
	registerProviders(cmd)
}

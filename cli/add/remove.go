package add

import "github.com/alecthomas/kingpin/v2"

func Register(app *kingpin.Application) {
	cmd := app.Command("add", "add an SCM provider")
	registerProvider(cmd)
}

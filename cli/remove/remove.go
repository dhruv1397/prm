package remove

import "github.com/alecthomas/kingpin/v2"

func Register(app *kingpin.Application) {
	cmd := app.Command("remove", "remove an SCM provider")
	registerProvider(cmd)
}

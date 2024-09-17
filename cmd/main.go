package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/prm/cli"
	"github.com/dhruv1397/prm/cli/add"
	"github.com/dhruv1397/prm/cli/list"
	"github.com/dhruv1397/prm/cli/purge"
	"github.com/dhruv1397/prm/cli/refresh"
	"github.com/dhruv1397/prm/cli/remove"
	"github.com/dhruv1397/prm/version"
)

const (
	application = "prm"
	description = "Pull Request Monitor (prm) is a CLI tool to list pull requests from different SCM providers."
)

func main() {
	args := cli.GetArguments()

	app := kingpin.New(application, description)
	list.Register(app)
	add.Register(app)
	remove.Register(app)
	refresh.Register(app)
	purge.Register(app)
	app.Version(version.Version.String())
	kingpin.MustParse(app.Parse(args))
}

package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/pr-monitor/cli"
	"github.com/dhruv1397/pr-monitor/cli/add"
	"github.com/dhruv1397/pr-monitor/cli/list"
	"github.com/dhruv1397/pr-monitor/cli/refresh"
	"github.com/dhruv1397/pr-monitor/cli/remove"
)

const (
	application = "prm"
	description = "Pull Request Monitor"
)

func main() {
	args := cli.GetArguments()

	app := kingpin.New(application, description)
	list.Register(app)
	add.Register(app)
	remove.Register(app)
	refresh.Register(app)
	kingpin.MustParse(app.Parse(args))
}

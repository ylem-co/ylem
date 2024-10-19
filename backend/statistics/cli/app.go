package cli

import (
	"ylem_statistics/cli/command/db"
	"ylem_statistics/cli/command/server"
	"ylem_statistics/cli/command/taskrun"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			db.Command,
			server.Command,
			taskrun.ResultListenerCommands,
		},
	}
}

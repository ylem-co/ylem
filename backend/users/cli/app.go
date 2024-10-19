package cli

import (
	"ylem_users/cli/command"
	"ylem_users/cli/command/encrypt"
	"ylem_users/cli/command/server"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			command.DbCommands,
			server.Command,
			encrypt.Command,
		},
	}
}

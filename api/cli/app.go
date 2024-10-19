package cli

import (
	"ylem_api/cli/command"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			command.ServerCommands,
			command.DbCommands,
		},
	}
}

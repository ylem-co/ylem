package cli

import (
	"ylem_integrations/cli/command/db"
	"ylem_integrations/cli/command/kafka"
	"ylem_integrations/cli/command/server"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			server.Command,
			kafka.Commands,
			db.DbCommands,
		},
	}
}

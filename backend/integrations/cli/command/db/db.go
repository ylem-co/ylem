package db

import "github.com/urfave/cli/v2"

var DbCommands = &cli.Command{
	Name:  "db",
	Usage: "Database management commands",
	Subcommands: []*cli.Command{
		MigrationsCommand,
	},
}

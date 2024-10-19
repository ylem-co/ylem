package cli

import (
	"ylem_pipelines/cli/command"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			command.ServerCommands,
			command.ScheduleGeneratorCommands,
			command.SchedulePublisherCommands,
			command.DbCommands,
			command.TriggerListenerCommands,
		},
	}
}

package cli

import (
	"ylem_taskrunner/cli/command"
	"ylem_taskrunner/cli/command/server"

	"github.com/urfave/cli/v2"
)

func NewApplication() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			command.KafkaCommands,
			server.Command,
			command.LoadBalancerCommands,
			command.TaskRunnerCommands,
		},
	}
}

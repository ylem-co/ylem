package command

import (
	"ylem_pipelines/services/schedule"

	"github.com/urfave/cli/v2"
)

var schedulePublisherStartHandler cli.ActionFunc = func(ctx *cli.Context) error {
	p, err := schedule.NewPublisher(ctx.Context)
	if err != nil {
		return err
	}

	return p.Start()
}

var SchedulePublisherStart = &cli.Command{
	Name:   "start",
	Usage:  "Start schedule publisher",
	Action: schedulePublisherStartHandler,
}

var SchedulePublisherCommands = &cli.Command{
	Name:  "schedulepub",
	Usage: "Schedule publisher commands",
	Subcommands: []*cli.Command{
		SchedulePublisherStart,
	},
}

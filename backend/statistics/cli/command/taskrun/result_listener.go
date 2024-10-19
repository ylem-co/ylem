package taskrun

import (
	"ylem_statistics/config"
	"ylem_statistics/services/taskrun"

	"github.com/lovoo/goka"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

var resultListenerStartHandler cli.ActionFunc = func(ctx *cli.Context) error {
	log.Info("Starting result listener...")

	l := taskrun.NewResultListener()

	cfg := config.Cfg().Kafka
	g := goka.DefineGroup(
		goka.Group(cfg.ConsumerGroupName),
		goka.Input(goka.Stream(cfg.TaskRunResultsTopic), new(messaging.MessageCodec), l.StoreResult),
	)

	p, err := goka.NewProcessor(cfg.BootstrapServers, g)
	if err != nil {
		log.Fatalf("error starting result listener: %v", err)
	} else {
		log.Info("Started.")
	}

	if err = p.Run(ctx.Context); err != nil {
		log.Fatalf("error running processor: %v", err)
	} else {
		log.Printf("Processor shutdown cleanly")
	}

	return nil
}

var ResultListenerStart = &cli.Command{
	Name:   "start",
	Usage:  "Start result listener",
	Action: resultListenerStartHandler,
}

var ResultListenerCommands = &cli.Command{
	Name:  "resultlistener",
	Usage: "Result listener commands",
	Subcommands: []*cli.Command{
		ResultListenerStart,
	},
}

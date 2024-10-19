package command

import (
	"ylem_pipelines/config"
	"ylem_pipelines/services/trigger/listener"

	"github.com/lovoo/goka"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

var triggerListenerStartHandler cli.ActionFunc = func(ctx *cli.Context) error {
	log.Info("Starting trigger listener...")

	l, err := listener.NewTriggerListener(ctx.Context)
	if err != nil {
		return err
	}

	cfg := config.Cfg().Kafka
	g := goka.DefineGroup(
		goka.Group(cfg.TriggerListenerConsumerGroupName),
		goka.Input(goka.Stream(cfg.TaskRunResultsTopic), new(messaging.MessageCodec), l.OnTaskRunResult),
		goka.Output(goka.Stream(cfg.TaskRunsTopic), new(messaging.MessageCodec)),
	)

	p, err := goka.NewProcessor(cfg.BootstrapServers, g)
	if err != nil {
		log.Fatalf("error starting trigger listener: %v", err)
	} else {
		log.Info("Started.")
	}
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx.Context); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	<-done

	return nil
}

var TriggerListenerStart = &cli.Command{
	Name:   "start",
	Usage:  "Start trigger listener",
	Action: triggerListenerStartHandler,
}

var TriggerListenerCommands = &cli.Command{
	Name:  "triggerlistener",
	Usage: "Trigger listener commands",
	Subcommands: []*cli.Command{
		TriggerListenerStart,
	},
}

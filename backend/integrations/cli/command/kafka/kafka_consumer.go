package kafka

import (
	"ylem_integrations/config"
	"ylem_integrations/services"

	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/lovoo/goka"
)

var kafkaConsumeHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Starting consumer...")

	cfg := config.Cfg().Kafka
	g := goka.DefineGroup(
		goka.Group(cfg.YlemIntegrationsConsumerGroupName),
		goka.Input(goka.Stream(cfg.NotificationTaskRunResultsTopic), new(messaging.MessageCodec), processMessage),
	)

	p, err := goka.NewProcessor(cfg.BootstrapServers, g)
	if err != nil {
		log.Fatalf("error starting consumer: %v", err)

		return err
	} else {
		log.Info("Started.")
	}

	return p.Run(c.Context)
}

var processMessage goka.ProcessCallback = func(ctx goka.Context, msg interface{}) {
	m := msg.(*messaging.Envelope)
	services.ProcessMessage(m.Msg)
	log.Debugf("Message processed. Key: %s", ctx.Key())
}

var Consume = &cli.Command{
	Name:   "consume",
	Usage:  "Consume related messages from kafka",
	Action: kafkaConsumeHandler,
}

var Commands = &cli.Command{
	Name:  "kafka",
	Usage: "Kafka related commands",
	Subcommands: []*cli.Command{
		Consume,
	},
}

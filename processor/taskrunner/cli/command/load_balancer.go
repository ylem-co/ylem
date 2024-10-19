package command

import (
	"context"
	"hash"
	"os"
	"syscall"
	"os/signal"
	"ylem_taskrunner/config"
	"ylem_taskrunner/internal/loadbalancer"

	"github.com/IBM/sarama"
	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

var loadBalancerStartHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Starting load balancer...")

	cfg := config.Cfg().Kafka
	g := goka.DefineGroup(
		goka.Group(cfg.LoadBalancerConsumerGroupName),
		goka.Input(goka.Stream(cfg.TaskRunsTopic), new(codec.Bytes), transmitMessage),
		goka.Output(goka.Stream(cfg.TaskRunsLoadBalancedTopic), new(codec.Bytes)),
	)

	p, err := goka.NewProcessor(cfg.BootstrapServers, g, goka.WithProducerBuilder(lbProducerBuilder))
	if err != nil {
		log.Fatalf("error starting load balancer: %v", err)
	} else {
		log.Info("Started.")
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		defer close(done)
		if err = p.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait // wait for SIGINT/SIGTERM
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)
	cancel() // gracefully stop processor
	<-done

	return nil
}

var lbCodec = &messaging.MessageCodec{}

var transmitMessage goka.ProcessCallback = func(ctx goka.Context, msg interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic while retransmitting a message, recovered and skipped\n", r)
		}
	}()

	cfg := config.Cfg().Kafka
	log.Debug("Message transferred. Key: " + ctx.Key())

	m := msg.([]byte)
	if len(m) > maxOutputLength {
		log.Warnf("Message is too big, skipping. Key: %s", ctx.Key())
		return
	}

	e, err := lbCodec.Decode(m)
	if err != nil {
		log.Warnf("Unable to decode message, skipping. Key: %s", ctx.Key())
	}

	envelope := e.(*messaging.Envelope)
	task := envelope.Msg.(messaging.TaskInterface)

	wfId := task.GetPipelineUuid()
	wfrId := task.GetPipelineRunUuid()
	orgId := task.GetOrganizationUuid()
	ctx.Emit(goka.Stream(cfg.TaskRunsLoadBalancedTopic), ctx.Key(), msg, goka.WithCtxEmitHeaders(goka.Headers{
		loadbalancer.HeaderPipelineId:    wfId[:],
		loadbalancer.HeaderPipelineRunId: wfrId[:],
		loadbalancer.HeaderOrganizationId: orgId[:],
	}))
}

func lbProducerBuilder(brokers []string, clientID string, hasher func() hash.Hash32) (goka.Producer, error) {
	config := goka.DefaultConfig()
	config.ClientID = clientID
	config.Producer.MaxMessageBytes = maxOutputLength * 2
	partitioner, err := loadbalancer.NewPartitioner()
	if err != nil {
		return nil, err
	}

	config.Producer.Partitioner = func(topic string) sarama.Partitioner {
		return partitioner
	}

	return goka.NewProducer(brokers, config)
}

var LoadBalancerStart = &cli.Command{
	Name:   "start",
	Usage:  "Start task load balancer",
	Action: loadBalancerStartHandler,
}

var LoadBalancerCommands = &cli.Command{
	Name:  "loadbalancer",
	Usage: "Load balancer commands",
	Subcommands: []*cli.Command{
		LoadBalancerStart,
	},
}

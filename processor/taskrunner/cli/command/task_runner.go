package command

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"os/signal"
	"ylem_taskrunner/config"
	"ylem_taskrunner/domain/runner"
	"ylem_taskrunner/services/redis"

	messaging "github.com/ylem-co/shared-messaging"
	"github.com/lovoo/goka"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const maxOutputLength = 50 * 1024 * 1024

var taskRunnerStartHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Starting task runner...")

	replaceGokaConfig()

	cfg := config.Cfg().Kafka
	g := goka.DefineGroup(
		goka.Group(cfg.TaskRunnerConsumerGroupName),
		goka.Input(goka.Stream(cfg.TaskRunsLoadBalancedTopic), new(messaging.MessageCodec), processMessage),
		goka.Output(goka.Stream(cfg.TaskRunResultsTopic), new(messaging.MessageCodec)),
		goka.Output(goka.Stream(cfg.QueryTaskRunResultsTopic), new(messaging.MessageCodec)),
		goka.Output(goka.Stream(cfg.NotificationTaskRunResultsTopic), new(messaging.MessageCodec)),
	)

	p, err := goka.NewProcessor(cfg.BootstrapServers, g)
	if err != nil {
		log.Fatalf("error starting task runner: %v", err)
	} else {
		log.Info("Started.")
	}
	ctx, cancel := context.WithCancel(context.Background())
	redis.Init(ctx)

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

var processMessage goka.ProcessCallback = func(ctx goka.Context, msg interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic while processing a message, recovered and skipped\n", r)
		}
	}()

	cfg := config.Cfg().Kafka
	m := msg.(*messaging.Envelope)
	trChan, callback := runner.RunTask(m.Msg, ctx.Context())

	for tr := range trChan {
		if tr == nil {
			continue
		}

		if len(tr.Output) > maxOutputLength {
			tr.IsSuccessful = false
			tr.Output = make([]byte, 0)
			tr.Errors = []messaging.TaskRunError{
				{
					Code: messaging.ErrorInternal,
					Message: fmt.Sprintf(
						"Output produced by the task is bigger than the maximum of %d bytes. Please reduce the data volume returned by the operation.",
						maxOutputLength,
					),
				},
			}
		}

		ctx.Emit(goka.Stream(cfg.TaskRunResultsTopic), ctx.Key(), messaging.NewEnvelope(tr))
		log.Debugf("Default bus message processed. Key: %s", ctx.Key())

		if callback != nil {
			callback(ctx, m.Msg, tr)
		}
	}
}

func replaceGokaConfig() {
	c := goka.DefaultConfig()
	c.Producer.MaxMessageBytes = maxOutputLength * 2

	goka.ReplaceGlobalConfig(c)
}

var TaskRunnerStart = &cli.Command{
	Name:   "start",
	Usage:  "Start listening to tasks and execute them",
	Action: taskRunnerStartHandler,
}

var TaskRunnerCommands = &cli.Command{
	Name:  "taskrunner",
	Usage: "Task runner commands",
	Subcommands: []*cli.Command{
		TaskRunnerStart,
	},
}

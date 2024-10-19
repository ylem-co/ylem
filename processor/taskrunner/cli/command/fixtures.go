package command

import (
	"fmt"
	"ylem_taskrunner/config"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

var fixtureLoadHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Loading fixtures...")
	cfg := config.Cfg().Kafka
	emitter, err := goka.NewEmitter(cfg.BootstrapServers, goka.Stream(cfg.TaskRunsTopic), new(messaging.MessageCodec))
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}

	defer emitter.Finish() //nolint:all

	for i := 0; i < 100; i++ {
		for _, msg := range getMessages() {
			err = emitter.EmitSync("", msg)
			if err != nil {
				log.Fatalf("error emitting message: %v", err)
			}

			fmt.Println("Loaded 1 message")
		}
	}

	log.Info("Done.")

	return nil
}

func getMessages() []*messaging.Envelope {
	taskUuid := uuid.New()
	pipelineUuid := uuid.New()
	return []*messaging.Envelope{
		messaging.NewEnvelope(&messaging.RunQueryTask{
			Task: messaging.Task{
				TaskUuid:         taskUuid,
				PipelineUuid:     pipelineUuid,
				CreatorUuid:      uuid.New(),
				OrganizationUuid: uuid.New(),
			},
			Query: "SELECT * FROM test_table;",
			Source: messaging.SQLIntegration{
				Uuid: uuid.New(),
				Type: "mysql",
				ConnectionType: "ssh",
				Host: []byte("127.0.0.1"),
				User: "foo",
				Password: []byte("bar"),
				Port: 3306,
				Database: "not_ylem_taskrunner",
				SshHost: []byte("localhost"),
				SshPort: 22,
				SshUser: "foo.bar",
			},
		}),
	}
}

var FixtureLoadCommand = &cli.Command{
	Name:   "load",
	Usage:  "Load fixtures into kafka topics",
	Action: fixtureLoadHandler,
}

var FixturesCommand = &cli.Command{
	Name:  "fixtures",
	Usage: "Kafka fixtures",
	Subcommands: []*cli.Command{
		FixtureLoadCommand,
	},
}

var KafkaCommands = &cli.Command{
	Name:  "kafka",
	Usage: "Kafka management commands",
	Subcommands: []*cli.Command{
		FixturesCommand,
	},
}

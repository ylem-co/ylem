package db

import (
	"time"
	"encoding/json"
	"math/rand"
	"ylem_statistics/domain/entity"
	"ylem_statistics/domain/entity/persister"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

var fixtureLoadHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Loading fixtures into database...")
	p := persister.Instance()

	rows := []entity.TaskRun{}

	execUuid := uuid.New()
	orgUuid := uuid.New()
	creatorUuid := uuid.New()
	pipelineUuid := uuid.MustParse("afc73ee8-ac45-4c12-983f-584038a93b0d")
	taskUuid1 := uuid.MustParse("3fd06c4f-54aa-4d53-a090-632512e4b45c")
	taskUuid2 := uuid.MustParse("21539f8c-11df-4b50-930d-78035e125955")

	// run #1, successful, executed between "2021-01-01 00:00:00" and "2021-01-01 00:00:10"
	runUuid := uuid.New()
	task1ExecutedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	task2ExecutedAt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:10")
	rows = appendRun(
		rows,
		true,
		task1ExecutedAt,
		task2ExecutedAt,
		pipelineUuid,
		runUuid,
		taskUuid1,
		taskUuid2,
		execUuid,
		orgUuid,
		creatorUuid,
	)

	// run #2, failed, executed between "2021-02-01 00:00:00" and "2021-02-01 00:00:15"
	runUuid = uuid.New()
	task1ExecutedAt, _ = time.Parse("2006-01-02 15:04:05", "2021-02-01 00:00:00")
	task2ExecutedAt, _ = time.Parse("2006-01-02 15:04:05", "2021-02-01 00:00:15")
	rows = appendRun(
		rows,
		false,
		task1ExecutedAt,
		task2ExecutedAt,
		pipelineUuid,
		runUuid,
		taskUuid1,
		taskUuid2,
		execUuid,
		orgUuid,
		creatorUuid,
	)

	// run #3, successful, executed between "2021-03-31 23:00:01" and "2021-04-01 00:00:12"
	runUuid = uuid.New()
	task1ExecutedAt, _ = time.Parse("2006-01-02 15:04:05", "2021-03-31 00:00:01")
	task2ExecutedAt, _ = time.Parse("2006-01-02 15:04:05", "2021-04-01 00:00:12")
	rows = appendRun(
		rows,
		true,
		task1ExecutedAt,
		task2ExecutedAt,
		pipelineUuid,
		runUuid,
		taskUuid1,
		taskUuid2,
		execUuid,
		orgUuid,
		creatorUuid,
	)

	for _, v := range rows {
		formatted, _ := json.MarshalIndent(v, "", "    ")
		err := p.CreateTaskRun(&v)
		log.Info(string(formatted))
		if err != nil {
			return err
		}
	}

	log.Info("Done.")

	return nil
}

func appendRun(
	rows []entity.TaskRun,
	isSuccessful bool,
	task1ExecutedAt time.Time,
	task2ExecutedAt time.Time,
	pipelineUuid uuid.UUID,
	runUuid uuid.UUID,
	taskUuid1 uuid.UUID,
	taskUuid2 uuid.UUID,
	execUuid uuid.UUID,
	orgUuid uuid.UUID,
	creatorUuid uuid.UUID,
) []entity.TaskRun {
	rows = append(
		rows,
		entity.TaskRun{
			Uuid:             uuid.New(),
			ExecutorUuid:     execUuid,
			OrganizationUuid: orgUuid,
			CreatorUuid:      creatorUuid,
			PipelineUuid:     pipelineUuid,
			PipelineRunUuid:  runUuid,
			TaskUuid:         taskUuid1,
			TaskType:         messaging.TaskTypeQuery,
			IsInitialTask:    true,
			IsFinalTask:      false,
			IsSuccessful:     true,
			IsFatalFailure:   false,
			ExecutedAt:       task1ExecutedAt,
			Duration:         uint32(rand.Intn(500) + 10),
		},
	)

	rows = append(
		rows,
		entity.TaskRun{
			Uuid:             uuid.New(),
			ExecutorUuid:     execUuid,
			OrganizationUuid: orgUuid,
			CreatorUuid:      creatorUuid,
			PipelineUuid:     pipelineUuid,
			PipelineRunUuid:  runUuid,
			TaskUuid:         taskUuid2,
			TaskType:         messaging.TaskTypeNotification,
			IsInitialTask:    false,
			IsFinalTask:      true,
			IsSuccessful:     isSuccessful,
			IsFatalFailure:   !isSuccessful,
			ExecutedAt:       task2ExecutedAt,
			Duration:         uint32(rand.Intn(700) + 50),
		},
	)

	return rows
}

var FixtureLoadCommand = &cli.Command{
	Name:   "load",
	Usage:  "Load fixtures into database",
	Action: fixtureLoadHandler,
}

var FixturesCommand = &cli.Command{
	Name:  "fixtures",
	Usage: "Database fixtures",
	Subcommands: []*cli.Command{
		FixtureLoadCommand,
	},
}

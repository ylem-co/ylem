package command

import (
	"database/sql"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/tasktrigger/types"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/pipeline/common"
	"ylem_pipelines/helpers"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

var fixtureLoadHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Loading fixtures...")
	db := helpers.DbConn()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	w := createPipeline(db)
	t1 := createTask1(db, w)
	createScheduleTriggerForTask(db, w, t1)
	t2 := createTask2(db, w)
	linkTasks(db, w, t1, t2)

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Info("Done.")

	return nil
}

func createPipeline(db *sql.DB) *pipeline.Pipeline {
	w := &pipeline.Pipeline{
		Uuid:    uuid.NewString(),
		Name:    "Test pipeline",
		Type:    common.PipelineTypeGeneric,
		Preview: make([]byte, 0),
	}
	id, err := pipeline.CreatePipeline(db, w)
	if err != nil {
		panic(err)
	}
	w.Id = int64(id)

	return w
}

func createTask1(db *sql.DB, w *pipeline.Pipeline) *task.Task {
	t := &task.Task{
		Uuid:         uuid.NewString(),
		PipelineId:   w.Id,
		PipelineUuid: w.Uuid,
		Type:         task.TaskTypeQuery,
		IsActive:     1,
	}
	reqTask := task.HttpApiNewTask{
		Name: "Test query",
		Type: task.TaskTypeQuery,
		Query: &task.HttpApiNewQuery{
			SQLQuery:   "SELECT * FROM sometable",
			SourceUuid: uuid.NewString(),
		},
	}
	id, _, err := task.CreateTaskWithImplementation(db, t, reqTask)
	if err != nil {
		panic(err)
	}
	t.Id = int64(id)

	return t
}

func createScheduleTriggerForTask(db *sql.DB, w *pipeline.Pipeline, t *task.Task) *tasktrigger.TaskTrigger {
	tt := &tasktrigger.TaskTrigger{
		Uuid:              uuid.NewString(),
		PipelineId:        w.Id,
		PipelineUuid:      w.Uuid,
		TriggeredTaskId:   t.Id,
		TriggeredTaskUuid: t.Uuid,
		TriggerType:       types.TriggerTypeSchedule,
		IsActive:          1,
		Schedule:          "* * * * *",
	}

	id, _, err := tasktrigger.CreateTaskTrigger(db, tt)
	if err != nil {
		panic(err)
	}

	tt.Id = int64(id)

	return tt
}

func createTask2(db *sql.DB, w *pipeline.Pipeline) *task.Task {
	t := &task.Task{
		Uuid:         uuid.NewString(),
		PipelineId:   w.Id,
		PipelineUuid: w.Uuid,
		Type:         task.TaskTypeCondition,
		IsActive:     1,
	}
	reqTask := task.HttpApiNewTask{
		Name: "Test condition",
		Type: task.TaskTypeCondition,
		Condition: &task.HttpApiNewCondition{
			Expression: "a > b",
		},
	}
	id, _, err := task.CreateTaskWithImplementation(db, t, reqTask)
	if err != nil {
		panic(err)
	}
	t.Id = int64(id)

	return t
}

func linkTasks(db *sql.DB, w *pipeline.Pipeline, triggerTask *task.Task, triggeredTask *task.Task) *tasktrigger.TaskTrigger {
	tt := &tasktrigger.TaskTrigger{
		Uuid:              uuid.NewString(),
		PipelineId:        w.Id,
		PipelineUuid:      w.Uuid,
		TriggerTaskId:     triggerTask.Id,
		TriggerTaskUuid:   triggerTask.Uuid,
		TriggeredTaskId:   triggeredTask.Id,
		TriggeredTaskUuid: triggeredTask.Uuid,
		TriggerType:       types.TriggerTypeOutput,
		IsActive:          1,
	}

	id, _, err := tasktrigger.CreateTaskTrigger(db, tt)
	if err != nil {
		panic(err)
	}

	tt.Id = int64(id)

	return tt
}

var FixtureLoadCommand = &cli.Command{
	Name:   "load",
	Usage:  "Load fixtures into database",
	Action: fixtureLoadHandler,
}

var FixturesCommand = &cli.Command{
	Name:  "fixtures",
	Usage: "Fixtures",
	Subcommands: []*cli.Command{
		FixtureLoadCommand,
	},
}

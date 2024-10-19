package messaging

import (
	"fmt"
	"reflect"
	"database/sql"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"

	messaging "github.com/ylem-co/shared-messaging"
)

type MergeTaskMessageFactory struct {
	db *sql.DB
}

func (f *MergeTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Merge)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Merge{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.MergeTask{
		Task:       task,
		FieldNames: impl.FieldNames,
	}

	inputCount, err := tasktrigger.GetInputCount(f.db, trc.Task.Id)
	msg.Task.Meta.InputCount = inputCount

	return messaging.NewEnvelope(msg), err
}

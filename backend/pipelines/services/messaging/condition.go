package messaging

import (
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"

	messaging "github.com/ylem-co/shared-messaging"
)

type ConditionTaskMessageFactory struct {
}

func (f *ConditionTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Condition)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Condition{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.CheckConditionTask{
		Task:       task,
		Expression: impl.Expression,
	}

	return messaging.NewEnvelope(msg), nil
}

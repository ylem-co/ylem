package messaging

import (
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"

	messaging "github.com/ylem-co/shared-messaging"
)

type AggregatorTaskMessageFactory struct {
}

func (f *AggregatorTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Aggregator)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Aggregator{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.AggregateDataTask{
		Task:         task,
		Expression:   impl.Expression,
		VariableName: impl.VariableName,
	}

	return messaging.NewEnvelope(msg), nil
}

package messaging

import (
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"

	messaging "github.com/ylem-co/shared-messaging"
)

type ExternalTriggerMessageFactory struct {
}

func (f *ExternalTriggerMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.ExternalTrigger)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.ExternalTrigger{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.ExternalTriggerTask{
		Task:  task,
		Input: trc.Input,
	}

	if string(trc.Input) == "{}" {
		msg.Input = []byte(impl.TestData)
	}

	return messaging.NewEnvelope(msg), nil
}

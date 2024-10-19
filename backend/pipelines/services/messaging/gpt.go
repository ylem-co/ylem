package messaging

import (
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"

	messaging "github.com/ylem-co/shared-messaging"
)

type GptMessageFactory struct {
}

func (f *GptMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Gpt)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Gpt{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.CallOpenapiGptTask{
		Task:   task,
		Prompt: impl.Prompt,
	}

	return messaging.NewEnvelope(msg), nil
}

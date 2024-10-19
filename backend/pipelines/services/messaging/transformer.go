package messaging

import (
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"

	messaging "github.com/ylem-co/shared-messaging"
)

type TransformerTaskMessageFactory struct {
}

func (f *TransformerTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Transformer)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Transformer{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.TransformDataTask{
		Task:                task,
		Type:                impl.Type,
		JsonQueryExpression: impl.JsonQueryExpression,
		Delimiter:           impl.Delimiter,
		CastToType:          impl.CastToType,
		DecodeFormat:        impl.DecodeFormat,
		EncodeFormat:        impl.EncodeFormat,
	}

	return messaging.NewEnvelope(msg), nil
}

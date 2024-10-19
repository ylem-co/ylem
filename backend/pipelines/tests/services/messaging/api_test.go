package tests

import (
	"reflect"
	"testing"
	"ylem_pipelines/app/task"
	ylem_messaging "ylem_pipelines/services/messaging"

	messaging "github.com/ylem-co/shared-messaging"
	"github.com/google/uuid"
)

func TestApiTaskMessageFactory_CreateMessage(t *testing.T) {
	type args struct {
		trc ylem_messaging.TaskRunContext
	}

	var mf *ylem_messaging.ApiTaskMessageFactory

	var wrongTaskImplementation task.Query
	wrongTask := &task.Task{Type: task.TaskTypeQuery, Implementation: wrongTaskImplementation}
	wrongTrc := ylem_messaging.TaskRunContext{Task: wrongTask}

	correctTaskImplementation := &task.ApiCall{Uuid: "20469af1-3c50-499f-8b76-7d2ae0c66575"}
	correctTask := &task.Task{
		Uuid: "20469af1-3c50-499f-8b76-7d2ae0c66575",
		PipelineUuid: "41cde0a5-e1fa-4d6e-b2f9-f5eed64c1584",
		OrganizationUuid: "8c09086c-96f2-41e4-ab90-13b5278e9ed8", 
		Type: task.TaskTypeAggregator, 
		Implementation: correctTaskImplementation,
	}
	correctTrc := ylem_messaging.TaskRunContext{Task: correctTask}
	taskUuid, _ := uuid.Parse(correctTask.Uuid)
	pipelineUuid, _ := uuid.Parse(correctTask.PipelineUuid)
	organizationUuid, _ := uuid.Parse(correctTask.OrganizationUuid)
	_ = &messaging.Envelope{
		Headers: map[string]string{"X-Message-Name": "tasks.api"},
		Msg: &messaging.CallApiTask{
			Task:         messaging.Task{
				TaskUuid: taskUuid,
				PipelineUuid: pipelineUuid,
				OrganizationUuid: organizationUuid,
			},
			Payload:   correctTaskImplementation.Payload,
			QueryString: correctTaskImplementation.QueryString,
		},
	}

	tests := []struct {
		name    string
		f       *ylem_messaging.ApiTaskMessageFactory
		args    args
		want    *messaging.Envelope
		wantErr bool
	}{
		{"False message type", mf, args{trc: wrongTrc}, nil, true},
		{"Incorrect destination", mf, args{trc: correctTrc}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ylem_messaging.ApiTaskMessageFactory{}
			got, err := f.CreateMessage(tt.args.trc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApiTaskMessageFactory.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiTaskMessageFactory.CreateMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

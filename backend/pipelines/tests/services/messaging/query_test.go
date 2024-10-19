package tests

import (
	"reflect"
	"testing"
	"ylem_pipelines/app/task"
	ylem_messaging "ylem_pipelines/services/messaging"

	messaging "github.com/ylem-co/shared-messaging"
)

func TestQueryTaskMessageFactory_CreateMessage(t *testing.T) {
	type args struct {
		trc ylem_messaging.TaskRunContext
	}

	var mf *ylem_messaging.QueryTaskMessageFactory

	var wrongTaskImplementation task.Notification
	wrongTask := &task.Task{Type: task.TaskTypeNotification, Implementation: wrongTaskImplementation}
	wrongTrc := ylem_messaging.TaskRunContext{Task: wrongTask}

	tests := []struct {
		name    string
		f       *ylem_messaging.QueryTaskMessageFactory
		args    args
		want    *messaging.Envelope
		wantErr bool
	}{
		{"False message type", mf, args{trc: wrongTrc}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ylem_messaging.QueryTaskMessageFactory{}
			got, err := f.CreateMessage(tt.args.trc)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryTaskMessageFactory.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryTaskMessageFactory.CreateMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

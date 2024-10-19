package tests

import (
	"ylem_users/entities"
	"testing"
)

func TestIsActionValid(t *testing.T) {
	type args struct {
		action string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{action: entities.ACTION_CREATE}, true},
		{"Valid 2", args{action: entities.ACTION_READ}, true},
		{"Valid 3", args{action: entities.ACTION_READ_LIST}, true},
		{"Valid 4", args{action: entities.ACTION_UPDATE}, true},
		{"Valid 5", args{action: entities.ACTION_RUN}, true},
		{"Valid 6", args{action: entities.ACTION_DELETE}, true},
		{"Invalid 1", args{action: "Some random action type"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.IsActionValid(tt.args.action); got != tt.want {
				t.Errorf("IsActionValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

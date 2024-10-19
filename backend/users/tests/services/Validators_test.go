package tests

import (
	"ylem_users/services"
	"testing"
)

func TestIsPhoneValid(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid 1", args{phone: "1(234)5678901x1234"}, true},
		{"Valid 2", args{phone: "(+351) 282 43 50 50"}, true},
		{"Valid 3", args{phone: "90191919908"}, true},
		{"Valid 4", args{phone: "555-8909"}, true},
		{"Valid 5", args{phone: "001 6867684"}, true},
		{"Valid 6", args{phone: "001 6867684x1"}, true},
		{"Valid 7", args{phone: "1 (234) 567-8901"}, true},
		{"Valid 8", args{phone: "1-234-567-8901 ext1234"}, true},
		{"Invalid 1", args{phone: "1-ab234-567-8901 ext1234"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsPhoneValid(tt.args.phone); got != tt.want {
				t.Errorf("IsPhoneValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPasswordValid(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Invalid 1", args{pwd: "12312313"}, false},
		{"Invalid 2", args{pwd: "123n123n"}, false},
		{"Invalid 3", args{pwd: "this"}, false},
		{"Valid 1", args{pwd: "thisIsPwd123&"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := services.IsPasswordValid(tt.args.pwd); got != tt.want {
				t.Errorf("IsPasswordValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

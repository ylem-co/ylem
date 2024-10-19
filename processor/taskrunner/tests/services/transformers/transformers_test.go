package tests

import (
	"reflect"
	"testing"
	"encoding/json"

	"ylem_taskrunner/services/transformers"
)

func TestSplitString(t *testing.T) {
	type args struct {
		value     string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Valid string splitting", args{value: "Something, else,is written,here", delimiter: ","}, []string{"Something", " else", "is written", "here"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformers.SplitString(tt.args.value, tt.args.delimiter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastStringToInteger(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"Valid transformation", args{value: "1234"}, 1234, false},
		{"Invalid transformation", args{value: "string"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformers.CastStringToInteger(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CastStringToInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CastStringToInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastFloatToInteger(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"Valid transformation", args{value: 1234.56}, 1234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformers.CastFloatToInteger(tt.args.value); got != tt.want {
				t.Errorf("CastFloatToInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastToStringType(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Valid transformation", args{value: 1234.56}, "1234.56"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformers.CastToStringType(tt.args.value); got != tt.want {
				t.Errorf("CastToStringType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeToCsv(t *testing.T) {
	testInput := []map[string]interface{}{
		{"id": 123, "name": "some string"},
	}
	encodedInput, _ := json.Marshal(testInput)

	type args struct {
		value        []byte
		delimiter    string
		columnsOrder []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Valid transformation", args{value: encodedInput, delimiter: ";", columnsOrder: []string{"name", "id"}}, []byte("name;id\nsome string;123\n"), false},
		{"Valid transformation. Wrong column in the order", args{value: encodedInput, delimiter: ";", columnsOrder: []string{"name", "uuid"}}, []byte("name;uuid\nsome string\n"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformers.EncodeToCsv(tt.args.value, tt.args.delimiter, tt.args.columnsOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToCsv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeToCsv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeToXml(t *testing.T) {
	testInput := []map[string]interface{}{
		{"id": 123, "name": "some string"},
	}
	encodedInput, _ := json.Marshal(testInput)

	type args struct {
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Valid transformation", args{value: encodedInput}, []byte("<xml>\n  <element>\n    <id>123</id>\n    <name>some string</name>\n  </element>\n</xml>"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformers.EncodeToXml(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeToXml() = %v, want %v", got, tt.want)
			}
		})
	}
}

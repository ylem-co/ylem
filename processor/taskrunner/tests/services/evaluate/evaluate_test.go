package tests

import (
	"context"
	"testing"
	"reflect"
	"encoding/json"

	"ylem_taskrunner/services/evaluate"
	hevaluate "ylem_taskrunner/helpers/evaluate"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestConditionWithContext(t *testing.T) {
	type args struct {
		ctx        context.Context
		expression string
		data       interface{}
	}

	tUuid, _ := uuid.FromBytes([]byte("f765be7b-c474-4fb7-8337-2481f1f2ebc6"))

	wrongTaskInput := struct{id int}{id: 123}
	ctxtWithWrongTaskInput := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    wrongTaskInput,
			EnvVars:      make(map[string]interface{}),
			PipelineUuid: tUuid,
		})

	var correctTaskInput interface{}
	_ = json.Unmarshal([]byte("[{id: 123}]"), &correctTaskInput)
	ctxtWithCorrectTaskInput := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    map[string]interface{}{"foo": 1, "bar": 2},
			EnvVars:      make(map[string]interface{}),
			PipelineUuid: tUuid,
		})

	strs := []map[string]interface{}{map[string]interface{}{"id": 1, "amount": 2.75, "created_at": "2020-01-01 00:00:00"}, map[string]interface{}{"id": 2, "amount": 34.58, "created_at": "2021-01-01 00:00:00"}, map[string]interface{}{"id": 3, "amount": 98.0, "created_at": "2022-01-01 00:00:00"}}
	arrayTaskInput := make([]interface{}, len(strs))
	for i, s := range strs {
	    arrayTaskInput[i] = s
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Invalid evaluation. Wrong TaskInput", args{ctx: ctxtWithWrongTaskInput, expression: "COUNT() > 0", data: make(map[string]interface{})}, false, true},
		{"False condition 1", args{ctx: ctxtWithCorrectTaskInput, expression: "COUNT() == 5", data: map[string]interface{}{"foo": 1, "bar": 2}}, false, false},
		{"True condition 1", args{ctx: ctxtWithCorrectTaskInput, expression: "bar === 2", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 2", args{ctx: ctxtWithCorrectTaskInput, expression: "bar >= foo", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 3", args{ctx: ctxtWithCorrectTaskInput, expression: "foo != null", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 4", args{ctx: ctxtWithCorrectTaskInput, expression: "foo !== \"bar\"", data: map[string]string{"foo": "foo", "bar": "bar"}}, true, false},
		{"True condition 5. Nil data", args{ctx: ctxtWithCorrectTaskInput, expression: "COUNT() > 0", data: nil}, true, false},
		{"True condition 6", args{ctx: ctxtWithCorrectTaskInput, expression: "bar == foo", data: map[string]string{"foo": "something", "bar": "something"}}, true, false},
		{"True condition 7", args{ctx: ctxtWithCorrectTaskInput, expression: "bar === foo", data: map[string]string{"foo": "something", "bar": "something"}}, true, false},
		{"True condition 8", args{ctx: ctxtWithCorrectTaskInput, expression: "\"2006-01-02\" < NOW()", data: map[string]string{"foo": "something", "bar": "something"}}, true, false},
		{"True condition 9", args{ctx: ctxtWithCorrectTaskInput, expression: "foo != bar", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 10", args{ctx: ctxtWithCorrectTaskInput, expression: "foo !== bar", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 11", args{ctx: ctxtWithCorrectTaskInput, expression: "foo <= bar", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 12", args{ctx: ctxtWithCorrectTaskInput, expression: "foo <== bar", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True condition 13", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(created_at) <== LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 14", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(created_at) <= LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 15", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(created_at) < LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 16", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(created_at) != LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 17", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(created_at) !== LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 18", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(created_at) == LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 19", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(created_at) === LAST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 20", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(created_at) >== FIRST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 21", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(created_at) >= FIRST(created_at)", data: arrayTaskInput}, true, false},
		{"True condition 22", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(created_at) > FIRST(created_at)", data: arrayTaskInput}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluate.ConditionWithContext(tt.args.ctx, tt.args.expression, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConditionWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConditionWithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateWithContext(t *testing.T) {
	type args struct {
		ctx        context.Context
		expression string
		data       interface{}
	}

	tUuid, _ := uuid.FromBytes([]byte("f765be7b-c474-4fb7-8337-2481f1f2ebc6"))

	wrongTaskInput := struct{id int}{id: 123}
	ctxtWithWrongTaskInput := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    wrongTaskInput,
			EnvVars:      make(map[string]interface{}),
			PipelineUuid: tUuid,
		})

	var correctTaskInput interface{}
	_ = json.Unmarshal([]byte("[{id: 123}]"), &correctTaskInput)
	ctxtWithCorrectTaskInput := context.WithValue(context.Background(), "ctx", hevaluate.Context{ //nolint:all
			TaskInput:    map[string]interface{}{"foo": 1, "bar": 2},
			EnvVars:      make(map[string]interface{}),
			PipelineUuid: tUuid,
		})

	strs := []map[string]interface{}{map[string]interface{}{"id": 1, "amount": 2.75, "created_at": "2020-01-01 00:00:00"}, map[string]interface{}{"id": 2, "amount": 34.58, "created_at": "2021-01-01 00:00:00"}, map[string]interface{}{"id": 3, "amount": 98.0, "created_at": "2022-01-01 00:00:00"}}
	arrayTaskInput := make([]interface{}, len(strs))
	for i, s := range strs {
	    arrayTaskInput[i] = s
	}

	strsOne := []map[string]interface{}{map[string]interface{}{"id": 1, "amount": 2.75}}
	arrayOneTaskInput := make([]interface{}, len(strsOne))
	for i, s := range strsOne {
	    arrayOneTaskInput[i] = s
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"True aggregation 1. Input is not used", args{ctx: ctxtWithWrongTaskInput, expression: "45", data: make(map[string]interface{})}, float64(45), false},
		{"False aggregation 1", args{ctx: ctxtWithCorrectTaskInput, expression: "COUNT()", data: map[string]interface{}{"foo": 1, "bar": 2}}, 1, false},
		{"True aggregation 2", args{ctx: ctxtWithCorrectTaskInput, expression: "bar + foo", data: map[string]interface{}{"foo": 1, "bar": 2}}, float64(3), false},
		{"True aggregation 3", args{ctx: ctxtWithCorrectTaskInput, expression: "bar >== foo", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True aggregation 4", args{ctx: ctxtWithCorrectTaskInput, expression: "foo != null", data: map[string]interface{}{"foo": 1, "bar": 2}}, true, false},
		{"True aggregation 5", args{ctx: ctxtWithCorrectTaskInput, expression: "foo != \"bar\"", data: map[string]string{"foo": "foo", "bar": "bar"}}, true, false},
		{"True aggregation 6", args{ctx: ctxtWithCorrectTaskInput, expression: "bar * foo - bar", data: map[string]interface{}{"foo": 3, "bar": 2}}, float64(4), false},
		{"True aggregation 7", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(bar / foo, 4, \"ceil\")", data: map[string]interface{}{"foo": 3, "bar": 2}}, float64(0.6667), false},
		{"True aggregation 8", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(bar / foo, 4, \"floor\")", data: map[string]interface{}{"foo": 3, "bar": 2}}, float64(0.6666), false},
		{"True aggregation 9", args{ctx: ctxtWithCorrectTaskInput, expression: "NEG(5)", data: nil}, decimal.NewFromFloat(float64(-5)), false},
		{"True aggregation 10", args{ctx: ctxtWithCorrectTaskInput, expression: "INT(5.1)", data: nil}, int64(5), false},
		{"True aggregation 11", args{ctx: ctxtWithCorrectTaskInput, expression: "STRING(5.1)", data: nil}, "5.1", false},
		{"True aggregation 12", args{ctx: ctxtWithCorrectTaskInput, expression: "SIGN(-5.1)", data: nil}, -1, false},
		{"True aggregation 13", args{ctx: ctxtWithCorrectTaskInput, expression: "SIGN(5.1)", data: nil}, 1, false},
		{"True aggregation 14", args{ctx: ctxtWithCorrectTaskInput, expression: "SIGN(0)", data: nil}, 0, false},
		{"True aggregation 15", args{ctx: ctxtWithCorrectTaskInput, expression: "ABS(-5.1)", data: nil}, float64(5.1), false},
		{"True aggregation 16", args{ctx: ctxtWithCorrectTaskInput, expression: "MIN(amount)", data: arrayTaskInput}, float64(2.75), false},
		{"True aggregation 17", args{ctx: ctxtWithCorrectTaskInput, expression: "MIN(amount)", data: arrayOneTaskInput}, float64(2.75), false},
		{"True aggregation 18", args{ctx: ctxtWithCorrectTaskInput, expression: "SUM(amount)", data: arrayTaskInput}, float64(135.33), false},
		{"True aggregation 19", args{ctx: ctxtWithCorrectTaskInput, expression: "MAX(amount)", data: arrayTaskInput}, float64(98.0), false},
		{"True aggregation 20", args{ctx: ctxtWithCorrectTaskInput, expression: "AVG(amount)", data: arrayTaskInput}, float64(45.11), false},
		{"True aggregation 21", args{ctx: ctxtWithCorrectTaskInput, expression: "FIRST(amount)", data: arrayTaskInput}, float64(2.75), false},
		{"True aggregation 22", args{ctx: ctxtWithCorrectTaskInput, expression: "LAST(amount)", data: arrayTaskInput}, float64(98.0), false},
		{"True aggregation 23", args{ctx: ctxtWithCorrectTaskInput, expression: "COUNT(amount)", data: arrayTaskInput}, int(3), false},
		{"True aggregation 24", args{ctx: ctxtWithCorrectTaskInput, expression: "MIN(created_at)", data: arrayTaskInput}, "2020-01-01 00:00:00", false},
		{"True aggregation 25", args{ctx: ctxtWithCorrectTaskInput, expression: "MAX(created_at)", data: arrayTaskInput}, "2022-01-01 00:00:00", false},
		{"True aggregation 25", args{ctx: ctxtWithCorrectTaskInput, expression: "INPUT()", data: arrayOneTaskInput}, "{\"bar\":2,\"foo\":1}", false},
		{"False aggregation 2. Input with args", args{ctx: ctxtWithCorrectTaskInput, expression: "INPUT(foo)", data: nil}, false, true},
		{"False aggregation 3. Bad ROUND args", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(bar < foo, 4)", data: map[string]interface{}{"foo": 3, "bar": 2}}, false, true},
		{"False aggregation 4. Bad ROUND args", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(\"string\", 4, \"ceil\")", data: map[string]interface{}{"foo": 3, "bar": 2}}, false, true},
		{"False aggregation 5. Bad ROUND args", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(2.45, \"string\", \"ceil\")", data: map[string]interface{}{"foo": 3, "bar": 2}}, false, true},
		{"False aggregation 6. Bad ROUND args", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(2.45, 1, 2)", data: map[string]interface{}{"foo": 3, "bar": 2}}, false, true},
		{"False aggregation 7. Bad ROUND args", args{ctx: ctxtWithCorrectTaskInput, expression: "ROUND(2.45, 1, \"wrong\")", data: map[string]interface{}{"foo": 3, "bar": 2}}, false, true},
		{"False aggregation 8. Bad ABS args", args{ctx: ctxtWithCorrectTaskInput, expression: "ABS()", data: nil}, false, true},
		{"False aggregation 9. Bad ABS args", args{ctx: ctxtWithCorrectTaskInput, expression: "ABS(\"string\")", data: nil}, false, true},
		{"False aggregation 10. Bad NEG args", args{ctx: ctxtWithCorrectTaskInput, expression: "NEG()", data: nil}, false, true},
		{"False aggregation 11. Bad NEG args", args{ctx: ctxtWithCorrectTaskInput, expression: "NEG(\"string\")", data: nil}, false, true},
		{"False aggregation 12. Bad STRING args", args{ctx: ctxtWithCorrectTaskInput, expression: "STRING()", data: nil}, false, true},
		{"False aggregation 13. Bad STRING args", args{ctx: ctxtWithCorrectTaskInput, expression: "STRING(\"string\")", data: nil}, false, true},
		{"False aggregation 14. Bad SIGN args", args{ctx: ctxtWithCorrectTaskInput, expression: "SIGN()", data: nil}, false, true},
		{"False aggregation 15. Bad SIGN args", args{ctx: ctxtWithCorrectTaskInput, expression: "SIGN(\"string\")", data: nil}, false, true},
		{"False aggregation 16. Bad INT args", args{ctx: ctxtWithCorrectTaskInput, expression: "INT()", data: nil}, false, true},
		{"False aggregation 17. Bad INT args", args{ctx: ctxtWithCorrectTaskInput, expression: "INT(\"string\")", data: nil}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluate.AggregateWithContext(tt.args.ctx, tt.args.expression, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateWithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

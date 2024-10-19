package runner

import (
	"time"
	"context"
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/redis"

	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
)

var errMalformedRequest error = errors.New("malformed request")

type InputItem struct {
	Input       []byte
	ColumnOrder []string
}

func MergeTaskRunner(t *messaging.MergeTask, ctx context.Context) *messaging.TaskRunResult {
	rc := redis.Instance()
	key := redisKey(t)
	ii := InputItem{
		Input:       t.Input,
		ColumnOrder: t.Meta.SqlQueryColumnOrder,
	}
	inputItemMarshalled, err := encodeInputItem(ii)
	if err != nil {
		log.Error(err)
		return malformedInput(t)
	}

	len, err := rc.LPush(ctx, key, inputItemMarshalled).Result()
	if err != nil {
		log.Error(err)
		return internalError(t)
	}

	_, err = rc.Expire(ctx, key, time.Minute*10).Result()
	if err != nil {
		log.Error(err)
		return internalError(t)
	}

	if len >= t.Meta.InputCount {
		inputItems, err := rc.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			log.Error(err)
			return internalError(t)
		}

		mergeFields := mergeFieldNames(t.FieldNames)
		result, columnOrder, err := doMerge(mergeFields, inputItems)
		if errors.Is(err, errMalformedRequest) {
			return malformedInput(t)
		}

		return mergeTaskRunResult(t, result, columnOrder)
	}

	return nil
}

func encodeInputItem(ii InputItem) (string, error) {
	encoded, err := json.Marshal(ii)

	return string(encoded), err
}

func decodeInputItem(inputItemMarshalled string) (InputItem, error) {
	ii := InputItem{}
	err := json.Unmarshal([]byte(inputItemMarshalled), &ii)

	return ii, err
}

func doMerge(mergeFields, inputItems []string) ([]map[string]interface{}, []string, error) {
	var result []map[string]interface{}
	columnOrder := make([]string, 0)
	for k, inputItemEncoded := range inputItems {
		inputItem, err := decodeInputItem(inputItemEncoded)
		if err != nil {
			return result, columnOrder, errMalformedRequest
		}

		rows, err := decodeInput(inputItem.Input)
		if err != nil {
			log.Error(err)
			return result, columnOrder, err
		}

		if k == 0 {
			result = rows
			columnOrder = inputItem.ColumnOrder
			continue
		}

		newResult := make([]map[string]interface{}, 0)
		mergedIdx := make(map[int]bool)
		for _, row1 := range result {
			merged := false
			for idx2, row2 := range rows {
				if rowsMatch(mergeFields, row1, row2) {
					newResult = append(newResult, mergeRows(row1, row2))
					merged = true
					mergedIdx[idx2] = true
				}
			}

			if !merged {
				newResult = append(newResult, row1)
			}
		}

		for idx2, row2 := range rows {
			if !mergedIdx[idx2] {
				newResult = append(newResult, row2)
			}
		}

		result = newResult

		columnOrder = mergeColumnOrder(columnOrder, inputItem.ColumnOrder)
	}

	return result, columnOrder, nil
}

func decodeInput(input []byte) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := json.Unmarshal(input, &result)
	if err != nil {
		var singleRow map[string]interface{}
		err = json.Unmarshal(input, &singleRow)
		if err != nil {
			log.Error(err)
			return result, errMalformedRequest
		}

		result = []map[string]interface{}{singleRow}
	}

	return result, nil
}

func mergeFieldNames(fieldNames string) []string {
	result := strings.Split(fieldNames, ",")
	for k, mfn := range result {
		result[k] = strings.TrimSpace(mfn)
	}

	if len(result) == 1 && result[0] == "" {
		return []string{}
	}

	return result
}

func rowsMatch(mergeFields []string, row1, row2 map[string]interface{}) bool {
	if len(mergeFields) == 0 {
		return true
	}

	for _, mf := range mergeFields {
		value1, ok1 := row1[mf]
		value2, ok2 := row2[mf]
		if !ok1 || !ok2 || value1 != value2 {
			return false
		}
	}

	return true
}

func mergeRows(rows ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, row := range rows {
		for k, v := range row {
			result[k] = v
		}
	}

	return result
}

func mergeColumnOrder(co1, co2 []string) []string {
	seen := make(map[string]bool)
	result := make([]string, len(co1))
	for idx, col := range co1 {
		seen[col] = true
		result[idx] = col
	}

	for _, col := range co2 {
		if seen[col] {
			continue
		}
		result = append(result, col)
	}

	return result
}

func mergeTaskRunResult(t *messaging.MergeTask, rows []map[string]interface{}, columnOrder []string) *messaging.TaskRunResult {
	return runMeasured(func() *messaging.TaskRunResult {
		tr := messaging.NewTaskRunResult(t.TaskUuid)

		tr.PipelineType = t.PipelineType
		tr.PipelineUuid = t.PipelineUuid
		tr.CreatorUuid = t.CreatorUuid
		tr.OrganizationUuid = t.OrganizationUuid
		tr.IsSuccessful = true
		tr.TaskType = messaging.TaskTypeMerge
		tr.PipelineRunUuid = t.PipelineRunUuid
		tr.TaskRunUuid = t.TaskRunUuid
		tr.IsInitialTask = t.IsInitialTask
		tr.IsFinalTask = t.IsFinalTask
		tr.Meta = t.Meta
		tr.Meta.SqlQueryColumnOrder = columnOrder

		var err error
		tr.Output, err = json.Marshal(rows)
		if err != nil {
			log.Error(err)
			tr.IsSuccessful = false
			tr.Errors = append(tr.Errors, messaging.TaskRunError{
				Code:     messaging.ErrorInternal,
				Severity: messaging.ErrorSeverityError,
				Message:  "Unable to produce output",
			})
		}

		return tr
	})
}

func internalError(t *messaging.MergeTask) *messaging.TaskRunResult {
	return &messaging.TaskRunResult{
		IsSuccessful:    false,
		PipelineRunUuid: t.PipelineRunUuid,
		TaskUuid:        t.TaskUuid,
		Errors: []messaging.TaskRunError{
			{
				Code:     messaging.ErrorInternal,
				Severity: messaging.ErrorSeverityError,
				Message:  "Internal error",
			},
		},
	}
}

func malformedInput(t *messaging.MergeTask) *messaging.TaskRunResult {
	return &messaging.TaskRunResult{
		IsSuccessful:    false,
		PipelineRunUuid: t.PipelineRunUuid,
		TaskUuid:        t.TaskUuid,
		Errors: []messaging.TaskRunError{
			{
				Code:     messaging.ErrorBadRequest,
				Severity: messaging.ErrorSeverityError,
				Message:  "Malformed input",
			},
		},
	}
}

func redisKey(t *messaging.MergeTask) string {
	cfg := config.Cfg().Redis

	return fmt.Sprintf(
		"%s%s.%s",
		cfg.KeyPrefix,
		"merge",
		t.PipelineRunUuid.String(),
	)
}

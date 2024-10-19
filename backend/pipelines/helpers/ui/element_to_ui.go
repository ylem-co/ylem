package ui

import (
	"encoding/json"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/tasktrigger/types"
)

type Elements struct {
	arr []map[string]interface{}
}

func (els *Elements) Add(el map[string]interface{}) {
	els.arr = append(els.arr, el)
}

func (els *Elements) MarshalJSON() ([]byte, error) {
	return json.Marshal(els.arr)
}

func NewElements() *Elements {
	return &Elements{
		arr: make([]map[string]interface{}, 0),
	}
}

func TaskToUi(t *task.Task, x, y int) map[string]interface{} {
	el := t.Implementation
	switch t.Type {
	case task.TaskTypeQuery:
		return queryToUi(t, el.(*task.Query), x, y)

	case task.TaskTypeNotification:
		return notificationToUi(t, el.(*task.Notification), x, y)
	}

	return nil
}

func TaskTriggerToUi(tt *tasktrigger.TaskTrigger) map[string]interface{} {
	var targetId interface{}
	targetId = 0
	if tt.TriggerTaskUuid != "" {
		targetId = tt.TriggerTaskUuid
	}
	result := map[string]interface{}{
		"id":           tt.Uuid,
		"trigger_type": tt.TriggerType,
		"source":       tt.TriggeredTaskUuid,
		"target":       targetId,
		"schedule":     tt.Schedule,
	}

	if tt.TriggerType != types.TriggerTypeSchedule {
		result["sourceHandle"] = nil
		result["targetHandle"] = nil
	}

	return result

}

func queryToUi(t *task.Task, el *task.Query, x, y int) map[string]interface{} {
	return map[string]interface{}{
		"id":       t.Uuid,
		"type":     task.TaskTypeQuery,
		"position": position(x, y),
		"data": map[string]interface{}{
			"name": t.Name,
		},
	}
}

func notificationToUi(t *task.Task, el *task.Notification, x, y int) map[string]interface{} {
	return map[string]interface{}{
		"id":       t.Uuid,
		"type":     task.TaskTypeNotification,
		"position": position(x, y),
		"data": map[string]interface{}{
			"name": t.Name,
		},
	}
}

func position(x, y int) map[string]int {
	return map[string]int{
		"x": x,
		"y": y,
	}
}

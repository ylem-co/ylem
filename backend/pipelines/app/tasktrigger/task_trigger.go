package tasktrigger

import (
	"regexp"
	"ylem_pipelines/app/tasktrigger/types"
)

type TaskTrigger struct {
	Id                int64  `json:"-"`
	Uuid              string `json:"uuid"`
	PipelineId        int64  `json:"-"`
	PipelineUuid      string `json:"pipeline_uuid"`
	TriggerType       string `json:"trigger_type"`
	Schedule          string `json:"schedule"`
	TriggerTaskUuid   string `json:"trigger_task_uuid"`
	TriggeredTaskUuid string `json:"triggered_task_uuid"`
	TriggerTaskId     int64  `json:"-"`
	TriggeredTaskId   int64  `json:"-"`
	IsActive          int8   `json:"-"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

func IsTriggerTypeSupported(Type string) bool {
	return map[string]bool{
		types.TriggerTypeSchedule:       true,
		types.TriggerTypeConditionTrue:  true,
		types.TriggerTypeConditionFalse: true,
		types.TriggerTypeOutput:         true,
	}[Type]
}

func IsScheduleValid(schedule string) bool {
	scheduleRegex, _ := regexp.Compile(`^(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|([\d\*]+(\/|-)\d+)|\d+|\*) ?){5,7})$`)
	return scheduleRegex.MatchString(schedule)
}

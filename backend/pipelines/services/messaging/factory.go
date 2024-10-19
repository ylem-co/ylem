package messaging

import (
	"context"
	"fmt"
	"time"
	"ylem_pipelines/app/task"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services/ylem_integrations"

	messaging "github.com/ylem-co/shared-messaging"
	"github.com/google/uuid"
)

type PipelineRunContext struct {
	ScheduledRunId   int64
	PipelineRunUuid  uuid.UUID
	OrganizationUuid string
	PipelineUuid     string
	PipelineType     string
	TaskRuns         []TaskRunContext
}

type TaskRunContext struct {
	PipelineRunContext PipelineRunContext
	PipelineRunUuid    uuid.UUID
	TaskRunUuid        uuid.UUID `json:"task_run_uuid"` // reserved for future use
	ExecuteAt          *time.Time
	Task               *task.Task
	IsInitialTask      bool `json:"is_initial_task"`
	IsFinalTask        bool `json:"is_final_task"`
	Input              []byte
	Meta               messaging.Meta
}

type MessageFactory interface {
	CreateMessage(TaskRunContext) (*messaging.Envelope, error)
}

type CompositeMessageFactory struct {
	queryTaskMessageFactory           *QueryTaskMessageFactory
	conditionTaskMessageFactory       *ConditionTaskMessageFactory
	aggregatorTaskMessageFactory      *AggregatorTaskMessageFactory
	transformerTaskMessageFactory     *TransformerTaskMessageFactory
	apiTaskMessageFactory             *ApiTaskMessageFactory
	notificationTaskMessageFactory    *NotificationTaskMessageFactory
	forEachTaskMessageFactory         *ForEachTaskMessageFactory
	mergeTaskMessageFactory           *MergeTaskMessageFactory
	filterTaskMessageFactory          *FilterTaskMessageFactory
	externalTriggerTaskMessageFactory *ExternalTriggerMessageFactory
	codeMessagingFactory              *CodeMessageFactory
	gptMessagingFactory               *GptMessageFactory
	processorTaskMessageFactory       *ProcessorTaskMessageFactory
	runPipelineTaskMessageFactory     *RunPipelineMessageFactory
}

func (f *CompositeMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	if trc.Task == nil {
		return nil, NewErrorNonRepeatable(fmt.Sprintf("message factory: pipeline run %s contains no task", trc.PipelineRunUuid.String()))
	}

	switch trc.Task.Type {
	case task.TaskTypeQuery:
		return f.queryTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeCondition:
		return f.conditionTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeAggregator:
		return f.aggregatorTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeTransformer:
		return f.transformerTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeApiCall:
		return f.apiTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeForEach:
		return f.forEachTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeNotification:
		return f.notificationTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeMerge:
		return f.mergeTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeFilter:
		return f.filterTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeExternalTrigger:
		return f.externalTriggerTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeCode:
		return f.codeMessagingFactory.CreateMessage(trc)

	case task.TaskTypeGpt:
		return f.gptMessagingFactory.CreateMessage(trc)

	case task.TaskTypeProcessor:
		return f.processorTaskMessageFactory.CreateMessage(trc)

	case task.TaskTypeRunPipeline:
		return f.runPipelineTaskMessageFactory.CreateMessage(trc)

	default:
		return nil, NewErrorNonRepeatable(fmt.Sprintf("task type %s is not supported", trc.Task.Type))
	}
}

func createTaskMessage(trc TaskRunContext) (messaging.Task, error) {
	t := messaging.Task{
		TaskRunUuid:     trc.TaskRunUuid,
		PipelineRunUuid: trc.PipelineRunUuid,
		PipelineType:    trc.PipelineRunContext.PipelineType,
		Input:           trc.Input,
		TaskName:        trc.Task.Name,
		IsInitialTask:   trc.IsInitialTask,
		IsFinalTask:     trc.IsFinalTask,
		Meta:            trc.Meta,
	}

	var err error
	tUid, err := uuid.Parse(trc.Task.Uuid)
	if err != nil {
		return t, fmt.Errorf("unable to parse task UUID: %s", err)
	}
	t.TaskUuid = tUid

	wfUid, err := uuid.Parse(trc.Task.PipelineUuid)
	if err != nil {
		return t, fmt.Errorf("unable to parse pipeline UUID: %s", err)
	}
	t.PipelineUuid = wfUid

	orgUid, err := uuid.Parse(trc.Task.OrganizationUuid)
	if err != nil {
		return t, fmt.Errorf("unable to parse organization UUID: %s", err)
	}
	t.OrganizationUuid = orgUid

	return t, nil
}

func createIntegrationMessage(dest ylem_integrations.Integration) (messaging.Integration, error) {
	d := messaging.Integration{
		Status:        dest.Status,
		Type:          dest.Type,
		Name:          dest.Name,
		Value:         dest.Value,
		UserUpdatedAt: dest.UserUpdatedAt,
	}

	dUid, err := uuid.Parse(dest.Uuid)
	if err != nil {
		return d, err
	}
	d.Uuid = dUid

	cUid, err := uuid.Parse(dest.CreatorUuid)
	if err != nil {
		return d, err
	}
	d.CreatorUuid = cUid

	orgUid, err := uuid.Parse(dest.OrganizationUuid)
	if err != nil {
		return d, err
	}
	d.OrganizationUuid = orgUid

	return d, nil
}

func NewCompositeMessageFactory(ctx context.Context) (MessageFactory, error) {
	var err error
	queryTaskMessageFactory, err := NewRunQueryTaskMessageFactory(ctx)
	if err != nil {
		return nil, err
	}

	apiTaskMessageFactory, err := NewApiTaskMessageFactory(ctx)
	if err != nil {
		return nil, err
	}

	notificationTaskMessageFactory, err := NewNotificationTaskMessageFactory(ctx)
	if err != nil {
		return nil, err
	}

	return &CompositeMessageFactory{
		queryTaskMessageFactory:        queryTaskMessageFactory,
		conditionTaskMessageFactory:    &ConditionTaskMessageFactory{},
		aggregatorTaskMessageFactory:   &AggregatorTaskMessageFactory{},
		transformerTaskMessageFactory:  &TransformerTaskMessageFactory{},
		forEachTaskMessageFactory:      &ForEachTaskMessageFactory{},
		apiTaskMessageFactory:          apiTaskMessageFactory,
		notificationTaskMessageFactory: notificationTaskMessageFactory,
		mergeTaskMessageFactory: &MergeTaskMessageFactory{
			db: helpers.DbConn(),
		},
		filterTaskMessageFactory:  &FilterTaskMessageFactory{},
		codeMessagingFactory:      &CodeMessageFactory{},
		gptMessagingFactory:       &GptMessageFactory{},
		processorTaskMessageFactory:    &ProcessorTaskMessageFactory{},
		runPipelineTaskMessageFactory:  &RunPipelineMessageFactory{},
	}, nil
}

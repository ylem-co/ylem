package task

import (
	"database/sql"
)

func getDBTable(Type string) string {
	return map[string]string{
		TaskTypeQuery:           "queries",
		TaskTypeCondition:       "conditions",
		TaskTypeAggregator:      "aggregators",
		TaskTypeTransformer:     "transformers",
		TaskTypeNotification:    "notifications",
		TaskTypeApiCall:         "api_calls",
		TaskTypeCode:            "codes",
		TaskTypeForEach:         "for_eaches",
		TaskTypeMerge:           "merges",
		TaskTypeFilter:          "filters",
		TaskTypeRunPipeline:     "run_pipelines",
		TaskTypeExternalTrigger: "external_triggers",
		TaskTypeProcessor:       "processors",
	}[Type]
}

func GetImplementation(db *sql.DB, task Task) (interface{}, error) {
	switch task.Type {
	case TaskTypeCondition:
		return GetCondition(db, task.ImplementationId)
	case TaskTypeAggregator:
		return GetAggregator(db, task.ImplementationId)
	case TaskTypeQuery:
		return GetQuery(db, task.ImplementationId)
	case TaskTypeNotification:
		return GetNotification(db, task.ImplementationId)
	case TaskTypeApiCall:
		return GetApiCall(db, task.ImplementationId)
	case TaskTypeTransformer:
		return GetTransformer(db, task.ImplementationId)
	case TaskTypeForEach:
		return GetForEach(db, task.ImplementationId)
	case TaskTypeMerge:
		return GetMerge(db, task.ImplementationId)
	case TaskTypeFilter:
		return GetFilter(db, task.ImplementationId)
	case TaskTypeRunPipeline:
		return GetRunPipeline(db, task.ImplementationId)
	case TaskTypeExternalTrigger:
		return GetExternalTrigger(db, task.ImplementationId)
	case TaskTypeCode:
		return GetCode(db, task.ImplementationId)
	case TaskTypeGpt:
		return GetGpt(db, task.ImplementationId)
	case TaskTypeProcessor:
		return GetProcessor(db, task.ImplementationId)
	default:
		return nil, nil
	}
}

func CreateImplementation(tx *sql.Tx, task HttpApiNewTask, wfUuid string) (int, error) {
	switch task.Type {
	case TaskTypeCondition:
		return CreateConditionTx(tx, task.Condition)
	case TaskTypeAggregator:
		return CreateAggregatorTx(tx, task.Aggregator, wfUuid)
	case TaskTypeQuery:
		return CreateQueryTx(tx, task.Query)
	case TaskTypeNotification:
		return CreateNotificationTx(tx, task.Notification)
	case TaskTypeApiCall:
		return CreateApiCallTx(tx, task.ApiCall)
	case TaskTypeForEach:
		return CreateForEachTx(tx, task.ForEach)
	case TaskTypeTransformer:
		return CreateTransformerTx(tx, task.Transformer)
	case TaskTypeMerge:
		return CreateMergeTx(tx, task.Merge)
	case TaskTypeFilter:
		return CreateFilterTx(tx, task.Filter)
	case TaskTypeRunPipeline:
		return CreateRunPipelineTx(tx, task.RunPipeline)
	case TaskTypeExternalTrigger:
		return CreateExternalTriggerTx(tx)
	case TaskTypeCode:
		return CreateCodeTx(tx, task.Code)
	case TaskTypeGpt:
		return CreateGptTx(tx, task.Gpt)
	case TaskTypeProcessor:
		return CreateProcessorTx(tx, task.Processor)
	default:
		return -1, nil
	}
}

func UpdateImplementation(tx *sql.Tx, task *Task, updatedTask HttpApiUpdatedTask) error {
	switch task.Type {
	case TaskTypeCondition:
		return UpdateConditionTx(tx, task.ImplementationId, updatedTask.Condition)
	case TaskTypeAggregator:
		return UpdateAggregatorTx(tx, task.ImplementationId, updatedTask.Aggregator)
	case TaskTypeQuery:
		return UpdateQueryTx(tx, task.ImplementationId, updatedTask.Query)
	case TaskTypeNotification:
		return UpdateNotificationTx(tx, task.ImplementationId, updatedTask.Notification)
	case TaskTypeApiCall:
		return UpdateApiCallTx(tx, task.ImplementationId, updatedTask.ApiCall)
	case TaskTypeTransformer:
		return UpdateTransformerTx(tx, task.ImplementationId, updatedTask.Transformer)
	case TaskTypeMerge:
		return UpdateMergeTx(tx, task.ImplementationId, updatedTask.Merge)
	case TaskTypeFilter:
		return UpdateFilterTx(tx, task.ImplementationId, updatedTask.Filter)
	case TaskTypeRunPipeline:
		return UpdateRunPipelineTx(tx, task.ImplementationId, updatedTask.RunPipeline)
	case TaskTypeExternalTrigger:
		return UpdateExternalTriggerTx(tx, task.ImplementationId, updatedTask.ExternalTrigger)
	case TaskTypeCode:
		return UpdateCodeTx(tx, task.ImplementationId, updatedTask.Code)
	case TaskTypeGpt:
		return UpdateGptTx(tx, task.ImplementationId, updatedTask.Gpt)
	case TaskTypeProcessor:
		return UpdateProcessorTx(tx, task.ImplementationId, updatedTask.Processor)
	default:
		return nil
	}
}

func CloneImplementation(tx *sql.Tx, taskType string, impl interface{}, wfUuid string) (int, error) {
	switch taskType {
	case TaskTypeCondition:
		return CloneConditionTx(tx, impl.(*Condition))
	case TaskTypeAggregator:
		return CloneAggregatorTx(tx, impl.(*Aggregator), wfUuid)
	case TaskTypeQuery:
		return CloneQueryTx(tx, impl.(*Query))
	case TaskTypeNotification:
		return CloneNotificationTx(tx, impl.(*Notification))
	case TaskTypeApiCall:
		return CloneApiCallTx(tx, impl.(*ApiCall))
	case TaskTypeForEach:
		return CloneForEachTx(tx, impl.(*ForEach))
	case TaskTypeTransformer:
		return CloneTransformerTx(tx, impl.(*Transformer))
	case TaskTypeMerge:
		return CloneMergeTx(tx, impl.(*Merge))
	case TaskTypeFilter:
		return CloneFilterTx(tx, impl.(*Filter))
	case TaskTypeRunPipeline:
		return CloneRunPipelineTx(tx, impl.(*RunPipeline))
	case TaskTypeExternalTrigger:
		return CloneExternalTriggerTx(tx)
	case TaskTypeCode:
		return CloneCodeTx(tx, impl.(*Code))
	case TaskTypeGpt:
		return CloneGptTx(tx, impl.(*Gpt))
	case TaskTypeProcessor:
		return CloneProcessorTx(tx, impl.(*Processor))
	default:
		return -1, nil
	}
}

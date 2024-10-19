package runner

import (
	"context"
	"ylem_taskrunner/config"
	"time"

	"github.com/ylem-co/shared-messaging/sources"
	"github.com/lovoo/goka"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
	"github.com/google/uuid"
)

type postTaskRunCallbackFn func(ctx goka.Context, t interface{}, tr *messaging.TaskRunResult)

func RunTask(task interface{}, ctx context.Context) (<-chan *messaging.TaskRunResult, postTaskRunCallbackFn) {
	result := make(chan *messaging.TaskRunResult)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic while running a task, recovered and skipped\n", r)
			}
		}()

		defer close(result)

		var sendResult = func(tr *messaging.TaskRunResult) {
			if tr != nil {
				tr.Uuid = uuid.New()
			}
			result <- tr
		}

		switch task := task.(type) {
		case *messaging.RunQueryTask:
			sendResult(
				RunQueryTaskRunner(task),
			)
		case *messaging.AggregateDataTask:
			sendResult(
				AggregateDataTaskRunner(task),
			)
		case *messaging.TransformDataTask:
			sendResult(
				TransformDataTaskRunner(task),
			)
		case *messaging.CheckConditionTask:
			sendResult(
				CheckConditionTaskRunner(task),
			)
		case *messaging.RunForEachTask:
			for _, tr := range RunForEachTaskRunner(task) {
				sendResult(
					tr,
				)
			}
		case *messaging.CallApiTask:
			sendResult(
				CallApiTaskRunner(task),
			)
		case *messaging.SendNotificationTask:
			sendResult(
				SendNotificationTaskRunner(task, ctx),
			)
		case *messaging.MergeTask:
			sendResult(
				MergeTaskRunner(task, ctx),
			)
		case *messaging.FilterTask:
			sendResult(
				FilterTaskRunner(task),
			)
		case *messaging.ExternalTriggerTask:
			sendResult(
				ExternalTriggerTaskRunner(task),
			)
		case *messaging.ProcessDataTask:
			sendResult(
				ProcessDataTaskRunner(task),
			)
		case *messaging.ExecuteCodeTask:
			sendResult(
				CodeTaskRunner(task),
			)
		case *messaging.CallOpenapiGptTask:
			sendResult(
				GptTaskRunner(task),
			)
		case *messaging.RunPipelineTask:
			sendResult(
				RunPipelineTaskRunner(task),
			)
		default:
			tr := &messaging.TaskRunResult{
				IsSuccessful: false,
				Errors: []messaging.TaskRunError{
					{
						Code:     messaging.ErrorUnknownTaskType,
						Severity: messaging.ErrorSeverityWarning,
						Message:  "Unknown task type",
					},
				},
			}

			sendResult(tr)
		}
	}()

	return result, getCallback(task)
}

func getCallback(task interface{}) postTaskRunCallbackFn {
	var cb postTaskRunCallbackFn
	switch task.(type) {
	case *messaging.RunQueryTask:
		cb = sendRunQueryTaskResultToSeparateTopic
	case *messaging.CallApiTask:
		cb = sendNotificationTaskResultToSeparateTopic
	case *messaging.SendNotificationTask:
		cb = sendNotificationTaskResultToSeparateTopic
	}

	return cb
}

func runMeasured(f func() *messaging.TaskRunResult) *messaging.TaskRunResult {
	start := time.Now()
	tr := f()
	tr.ExecutedAt = time.Now()
	tr.Duration = time.Since(start) / time.Millisecond

	return tr
}

func runMeasuredMultiOutput(f func() []*messaging.TaskRunResult) []*messaging.TaskRunResult {
	start := time.Now()
	trs := f()
	for _, tr := range trs {
		tr.ExecutedAt = time.Now()
		tr.Duration = time.Since(start) / time.Millisecond
	}

	return trs
}

func sendRunQueryTaskResultToSeparateTopic(ctx goka.Context, t interface{}, tr *messaging.TaskRunResult) {
	rqt := t.(*messaging.RunQueryTask)

	if (!tr.IsSuccessful && rqt.Source.Status == messaging.SQLIntegrationStatusOffline) ||
		(tr.IsSuccessful && rqt.Source.Status == messaging.SQLIntegrationStatusOnline) {
		return
	}

	m := sources.SourceStatusToggled{Uuid: rqt.Source.Uuid.String()}
	cfg := config.Cfg().Kafka
	ctx.Emit(goka.Stream(cfg.QueryTaskRunResultsTopic), ctx.Key(), messaging.NewEnvelope(&m))
	log.Debugf(`"%s" message processed. Key: %s`, sources.SourceStatusToggledMessageName, ctx.Key())
}

func sendNotificationTaskResultToSeparateTopic(ctx goka.Context, t interface{}, tr *messaging.TaskRunResult) {
	if tr.IsSuccessful {
		return
	}

	cfg := config.Cfg().Kafka
	ctx.Emit(goka.Stream(cfg.NotificationTaskRunResultsTopic), ctx.Key(), messaging.NewEnvelope(tr))
	log.Debugf(`"%s" message processed. Key: %s`, messaging.TaskSendNotificationMessageName, ctx.Key())
}

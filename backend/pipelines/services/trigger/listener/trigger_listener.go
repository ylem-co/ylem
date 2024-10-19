package listener

import (
	"time"
	"context"
	"database/sql"
	"encoding/json"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/task/result"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/tasktrigger/types"
	"ylem_pipelines/app/pipeline/run"
	"ylem_pipelines/config"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services/kafka"
	msgsrv "ylem_pipelines/services/messaging"
	"ylem_pipelines/services/provider"
	"ylem_pipelines/services/schedule"
	"ylem_pipelines/services/taskrunner"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

type TriggerListener struct {
	topic              string
	db                 *sql.DB
	ctx                context.Context
	taskProvider       provider.TaskProvider
	msgFactory         msgsrv.MessageFactory
}

func (l *TriggerListener) OnTaskRunResult(ctx goka.Context, envelope interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic in trigger listener, recovered and skipped\n", r)
		}
	}()
	evl, ok := envelope.(*messaging.Envelope)
	if !ok {
		log.Trace("Unknown envelope type, skipping.")
		return
	}

	log.Trace("Trigger listener: got message", evl.Msg)

	trr, ok := evl.Msg.(*messaging.TaskRunResult)
	if !ok {
		log.Trace("Unknown message type, skipping.")
		return
	}

	isUpdated, err := l.updateStoredTaskRunResult(trr)
	if err != nil {
		ctx.Fail(err)
	}

	if !trr.IsSuccessful {
		log.Debugf("Task %s execution failed: %s", trr.TaskUuid.String(), trr.Errors[0].Message)
		return
	}

	if trr.TaskType == task.TaskTypeRunPipeline {
		tx, err := l.db.Begin()
		if err != nil {
			log.Error(err)
		} else {
			t, err := task.GetTaskByUuid(l.db, trr.TaskUuid.String())
			if err != nil {
				log.Debugf("Task %s not found, skipping", trr.TaskUuid)
				_ = tx.Rollback()
			} else {
				err = taskrunner.InitiatePipelineRun(tx, msgsrv.TaskRunContext{
					Task: t,
				}, []byte{}, uuid.Nil)

				if err != nil {
					_ = tx.Rollback()
					log.Error(err)
				} else {
					err = tx.Commit()
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

	if trr.IsFinalTask {
		log.Debugf("Task result %s is final, won't create new tasks", trr.TaskUuid.String())
		return
	}

	msgPipelineRunConfig := trr.Meta.PipelineRunConfig
	wrConfig := run.PipelineRunConfig{
		TaskIds: run.IdList{
			Type: msgPipelineRunConfig.TaskIds.Type,
			Ids:  msgPipelineRunConfig.TaskIds.Ids,
		},
		TaskTriggerIds: run.IdList{
			Type: msgPipelineRunConfig.TaskTriggerIds.Type,
			Ids:  msgPipelineRunConfig.TaskTriggerIds.Ids,
		},
	}

	var triggerType string
	var nextInput []byte
	cr := &messaging.ConditionResult{}
	switch trr.TaskType {
	case task.TaskTypeCondition:
		err = json.Unmarshal(trr.Output, cr)
		if err != nil {
			log.Error(err)
			return
		}

		if cr.Result {
			triggerType = types.TriggerTypeConditionTrue
		} else {
			triggerType = types.TriggerTypeConditionFalse
		}
		nextInput = cr.OriginalInput
	default:
		nextInput = trr.Output
		triggerType = types.TriggerTypeOutput
	}

	taskIds, err := tasktrigger.GetTriggeredTaskIds(l.db, trr.TaskUuid, triggerType, wrConfig)
	if err != nil {
		ctx.Fail(err)
	}

	if len(taskIds) == 0 {
		log.Debugf("No new tasks triggered by task %s", trr.TaskUuid.String())
		return
	}

	log.Debugf("Found %d triggered tasks", len(taskIds))

	for _, id := range taskIds {
		t, isFinal, err := l.taskProvider.GetTask(id, wrConfig)
		if err != nil {
			log.Debugf("Task %d not found, skipping", id)
			continue
		}

		key, msg, h, err := l.createNextTaskMessage(trr, t, isFinal, nextInput, isUpdated)
		if err != nil {
			if _, ok := err.(msgsrv.ErrorNonRepeatable); ok {
				log.Errorf("non-repeatable error, scheduled run will be deleted: %s", err)
				continue
			}

			if _, ok := err.(msgsrv.ErrorRepeatable); ok {
				log.Errorf("repeatable error, scheduled run will be retried: %s", err)
				ctx.Fail(err)
				return
			}

			ctx.Fail(err)
			return
		}

		ctx.Emit(
			goka.Stream(l.topic),
			key,
			msg,
			goka.WithCtxEmitHeaders(h),
		)
	}
}

func (l *TriggerListener) createNextTaskMessage(prevTrr *messaging.TaskRunResult, t *task.Task, isFinal bool, nextInput []byte, storeResult bool) (string, *messaging.Envelope, goka.Headers, error) {
	var h goka.Headers = make(goka.Headers)

	log.Debugf("Task %d is triggered. Scheduling it for immediate execution.", t.Id)

	if len(nextInput) > kafka.MaxTaskInputLength {
		errMsg := "Task input is too big"
		tUid, err := uuid.Parse(t.Uuid)
		if err != nil {
			return "", nil, h, msgsrv.NewErrorNonRepeatable(err.Error())
		}

		immediateResult := &messaging.TaskRunResult{
			PipelineRunUuid: prevTrr.PipelineRunUuid,
			TaskRunUuid:     prevTrr.TaskRunUuid,
			TaskUuid:        tUid,
			IsSuccessful:    false,
			Errors: []messaging.TaskRunError{
				{
					Code:     messaging.ErrorBadRequest,
					Severity: messaging.ErrorSeverityError,
					Message:  errMsg,
				},
			},
		}

		_, err = l.updateStoredTaskRunResult(immediateResult)
		if err != nil {
			return "", nil, h, msgsrv.NewErrorNonRepeatable(err.Error())
		}

		return "", nil, h, msgsrv.NewErrorNonRepeatable(errMsg)
	}

	now := time.Now()
	trc := msgsrv.TaskRunContext{
		PipelineRunContext: msgsrv.PipelineRunContext{
			PipelineType: prevTrr.PipelineType,
		},
		PipelineRunUuid: prevTrr.PipelineRunUuid,
		TaskRunUuid:     uuid.Nil,
		ExecuteAt:       &now,
		Task:            t,
		IsInitialTask:   false,
		IsFinalTask:     isFinal,
		Input:           nextInput,
		Meta:            prevTrr.Meta,
	}

	if storeResult {
		taskUuid, _ := uuid.Parse(t.Uuid)
		_, err := result.CreatePendingTaskRunResult(l.db, t.Id, taskUuid, trc.TaskRunUuid, trc.PipelineRunUuid)
		if err != nil {
			return "", nil, nil, err
		}
	}

	msg, err := l.msgFactory.CreateMessage(trc)
	if err != nil {
		return "", nil, nil, err
	}

	h = schedule.GetMessageHeaders(l.topic, trc)

	return trc.PipelineRunUuid.String(), msg, h, nil
}

func (l *TriggerListener) updateStoredTaskRunResult(trr *messaging.TaskRunResult) (bool, error) {
	return result.UpdateTaskRunResult(l.db, trr.TaskUuid, &result.TaskRunResult{
		TaskRunUuid:     trr.TaskRunUuid,
		PipelineRunUuid: trr.PipelineRunUuid,
		IsSuccessful:    trr.IsSuccessful,
		Output:          trr.Output,
		Errors:          taskRunResultErrors(trr),
	})
}

func taskRunResultErrors(trr *messaging.TaskRunResult) []result.TaskRunError {
	var errors []result.TaskRunError
	for _, e := range trr.Errors {
		errors = append(errors, result.TaskRunError{
			Code:     e.Code,
			Severity: e.Severity,
			Message:  e.Message,
		})
	}

	return errors
}

func NewTriggerListener(ctx context.Context) (*TriggerListener, error) {
	cfg := config.Cfg().Kafka
	db := helpers.DbConn()

	tp := &provider.DbTaskProvider{
		Db: db,
	}

	ctp := &provider.CachingTaskProvider{
		InnerProvider: tp,
		Ctx:           ctx,
	}

	err := ctp.Init()
	if err != nil {
		return nil, err
	}

	mf, err := msgsrv.NewCompositeMessageFactory(ctx)
	if err != nil {
		return nil, err
	}

	l := &TriggerListener{
		topic:              cfg.TaskRunsTopic,
		db:                 db,
		taskProvider:       ctp,
		msgFactory:         mf,
		ctx:                ctx,
	}

	log.Trace("Trigger listener initialized")

	return l, nil
}

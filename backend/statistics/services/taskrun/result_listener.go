package taskrun

import (
	"encoding/json"
	"ylem_statistics/domain/entity"
	"ylem_statistics/domain/entity/persister"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

type ResultListener struct {
	ep persister.EntityPersister
}

func (l *ResultListener) StoreResult(ctx goka.Context, envelope interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic in result listener, recovered and skipped\n", r)
		}
	}()

	evl, ok := envelope.(*messaging.Envelope)
	if !ok {
		log.Trace("Unknown envelope type, skipping.")
		return
	}

	log.Trace("Got message", evl.Msg)

	trr, ok := evl.Msg.(*messaging.TaskRunResult)
	if !ok || trr.PipelineRunUuid == uuid.Nil {
		log.Trace("Unknown message type, skipping.")
		return
	}

	tr := createFromMessage(trr)
	err := l.ep.CreateTaskRun(&tr)
	if err != nil {
		ctx.Fail(err)
	}
}

func createFromMessage(trr *messaging.TaskRunResult) entity.TaskRun {
	run := entity.TaskRun{
		Uuid:             trr.Uuid,
		ExecutorUuid:     uuid.Nil,
		OrganizationUuid: trr.OrganizationUuid,
		CreatorUuid:      trr.CreatorUuid,
		PipelineUuid:     trr.PipelineUuid,
		PipelineRunUuid:  trr.PipelineRunUuid,
		TaskUuid:         trr.TaskUuid,
		PipelineType:     (trr.PipelineType),
		TaskType:         (trr.TaskType),
		IsInitialTask:    trr.IsInitialTask,
		IsFinalTask:      trr.IsFinalTask,
		IsSuccessful:     trr.IsSuccessful,
		IsFatalFailure:   !trr.IsSuccessful,
		ExecutedAt:       trr.ExecutedAt,
		Duration:         uint32(trr.Duration),
	}

	if trr.TaskType == messaging.TaskTypeAggregator && run.PipelineType == messaging.PipelineTypeMetric && run.IsSuccessful {
		metricOutput := []map[string]float64{}
		err := json.Unmarshal(trr.Output, &metricOutput)
		if err != nil {
			return run
		}

		if len(metricOutput) == 0 {
			return run
		}

		row := metricOutput[0]
		value, ok := row["value"]
		if !ok {
			return run
		}

		run.IsMetricValueSet = true
		run.MetricValue = value
	}

	if run.PipelineType != messaging.PipelineTypeMetric {
		if !run.IsSuccessful && len(trr.Errors) > 0 {
			run.Output, _ = json.Marshal(trr.Errors)
		} else {
			run.Output = trr.Output
		}
	}

	return run
}

func NewResultListener() *ResultListener {
	return &ResultListener{
		ep: persister.Instance(),
	}
}

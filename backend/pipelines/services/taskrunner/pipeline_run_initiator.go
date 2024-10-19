package taskrunner

import (
	"fmt"
	"time"
	"database/sql"
	"encoding/json"
	"ylem_pipelines/app/envvariable"
	"ylem_pipelines/app/schedule"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/task/result"
	"ylem_pipelines/app/pipeline"
	msgsrv "ylem_pipelines/services/messaging"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func InitiatePipelineRun(tx *sql.Tx, trc msgsrv.TaskRunContext, prevOutputBytes []byte, prewPipelineRunUuid uuid.UUID) error {
	sysVars, err := envvariable.GetEnvVariablesByOrganizationUuidTx(tx, trc.Task.OrganizationUuid)
	if err != nil {
		return err
	}

	envVars, err := outputToEnvVars(prevOutputBytes)
	if err != nil {
		return err
	}

	for _, sv := range sysVars.Items {
		if _, ok := envVars[sv.Name]; ok {
			continue
		}

		envVars[sv.Name] = sv.Value
	}

	wf, err := pipeline.GetPipelineByUuidTx(tx, trc.Task.Implementation.((*task.RunPipeline)).PipelineUuid)
	if err != nil {
		return err
	}

	if wf == nil {
		return nil
	}

	executeAt := time.Now()
	wrUuid := uuid.New()
	sr := schedule.ScheduledRun{
		PipelineRunUuid: wrUuid,
		PipelineId:      wf.Id,
		Input:           make([]byte, 0),
		EnvVars:         envVars,
		ExecuteAt:       &executeAt,
	}

	err = schedule.AddScheduledRunsTx(tx, []schedule.ScheduledRun{sr})

	var output []byte
	if err != nil {
		output = []byte("Pipeline run failed")
	} else {
		output = []byte("Pipeline run initiated")
	}

	if prewPipelineRunUuid != uuid.Nil {
		tUid, err := uuid.Parse(trc.Task.Uuid)
		if err != nil {
			return fmt.Errorf("Pipeline initiation failed: %s", err)
		}

		_, err = result.UpdateTaskRunResultTx(tx, tUid, &result.TaskRunResult{
			PipelineRunUuid: prewPipelineRunUuid,
			IsSuccessful:    err == nil,
			Output:          output,
		})
		if err != nil {
			return err
		}
	}

	return err
}

func outputToEnvVars(output []byte) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if len(output) == 0 {
		return result, nil
	}

	var prevOutput []map[string]interface{}
	err := json.Unmarshal(output, &prevOutput)
	if err != nil {
		logrus.Error(err)
		return result, err
	}

	if len(prevOutput) == 0 {
		return result, nil
	}

	for key, value := range prevOutput[0] {
		result[key] = value
	}

	return result, nil
}

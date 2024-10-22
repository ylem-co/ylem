package trial

import (
	"database/sql"

	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/pipelinetemplate"

	log "github.com/sirupsen/logrus"
)

type DemoPipeline struct {
	Uuid string
	Name string
}

func CreateTrialPipelines(db *sql.DB, organizationUuid string, sourceUuid string, destinationUuid string, userUuid string) error {
	demoPipelines := []DemoPipeline {
		{"b518cb04-47a7-49fe-aa22-82106a080306", "Demo Pipeline: Basic Streaming from DB to API"},
		{"4fbe4c86-4ff3-41e8-aea5-062afa9655af", "Demo Pipeline: Generic Multiconditional Streaming"},
		{"044617b1-4f0f-4966-84af-d108006e1641", "Demo Pipeline: Processing Dataset with jq"},
		{"06927207-bb5a-4281-a81a-aadf2964004e", "Demo Pipeline: Processing Dataset with Python"},
		{"3d491dd4-b33c-435c-a642-6c1bbf8f7b74", "Demo Pipeline: Messaging Patterns. Message Router"},
		{"373ddfc2-efd8-444e-9438-804db0c6e749", "Demo Metric: ARR (Annual recurrent revenue)"},
		{"901ea1f5-c9e8-41ec-8adb-96030e2a39ff", "Demo Metric: OKR. HR. Monthly hires"},
		{"82f3d55b-6fd5-4397-b768-05be19a00487", "Demo Metric: OKR. Customer success. Number of sign ups"},
		{"ff4ae672-17ce-46c2-aab9-992e9067a61a", "Demo Metric: Average invoice amount"},
		{"835d2e46-be77-4707-aba3-0444276662ae", "Demo Metric: Number of complete units"},
	}
	for _, p := range demoPipelines {
		err := createTrialPipeline(db, p.Uuid, p.Name, organizationUuid, sourceUuid, destinationUuid, userUuid)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTrialPipeline(db *sql.DB, pipelineUuid string, pipelineName string, organizationUuid string, sourceUuid string, destinationUuid string, userUuid string) error {
	tplWf, tasks, tts, err := loadPipelineData(db, pipelineUuid)
	if err != nil {
		log.Error(err)
		return err
	}

	if tplWf == nil {
		log.Error(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	newWfUuid, err := pipelinetemplate.ClonePipelineTx(tx, organizationUuid, userUuid, tplWf, tasks, tts, 0, nil)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	newWf, err := pipeline.GetPipelineByUuidTx(tx, newWfUuid)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	newWf.IsTemplate = 0
	newWf.Name = pipelineName
	err = pipeline.UpdatePipelineTx(tx, newWf, false)
	if err != nil {
		log.Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}


	tx, err = db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	_, tasks, _, err = loadPipelineData(db, newWfUuid)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)
		return err
	}

	var query *task.Query
	var uQuery *task.HttpApiUpdatedQuery
	for _, qTask := range tasks.Items {
	    if qTask.Type == task.TaskTypeQuery {
	    	query, _ = task.GetQuery(db, qTask.Id)
	    	uQuery.SourceUuid = sourceUuid
	    	uQuery.SQLQuery = query.SQLQuery
	    	_ = task.UpdateQueryTx(tx, qTask.Id, uQuery)
	    }
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func loadPipelineData(db *sql.DB, pipelineUuid string) (*pipeline.Pipeline, *task.Tasks, *tasktrigger.TaskTriggers, error) {
	wf, err := pipeline.GetPipelineByUuid(db, pipelineUuid)
	if err != nil {
		return nil, nil, nil, err
	}

	if wf == nil {
		return nil, nil, nil, nil
	}

	tasks, err := task.GetTasksByPipelineId(db, wf.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	tts, err := tasktrigger.GetTaskTriggersByPipelineId(db, wf.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	return wf, tasks, tts, nil
}

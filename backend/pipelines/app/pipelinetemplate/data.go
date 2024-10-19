package pipelinetemplate

import (
	"errors"
	"strings"
	"database/sql"
	"ylem_pipelines/app/folder"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/pipeline"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	log "github.com/sirupsen/logrus"
)

func ClonePipelineTx(tx *sql.Tx, orgUuid, creatorUuid string, wf *pipeline.Pipeline, tasks *task.Tasks, triggers *tasktrigger.TaskTriggers, isTemplate int8, folder *folder.Folder) (string, error) {
	if wf.IsActive == 0 {
		return "", errors.New("unable to clone an inactive pipeline")
	}

	newWf := *wf
	newWf.Id = 0
	newWf.IsTemplate = isTemplate
	newWf.Schedule = ""
	newWf.OrganizationUuid = orgUuid
	newWf.CreatorUuid = creatorUuid

	if folder != nil {
		newWf.FolderId = folder.Id
		newWf.FolderUuid = folder.Uuid
	} else {
		newWf.FolderId = 0
		newWf.FolderUuid = ""
	}

	newWfId, err := pipeline.CreatePipelineTx(tx, &newWf)
	newWf.Id = int64(newWfId)
	if err != nil {
		_ = tx.Rollback()
		return uuid.Nil.String(), err
	}

	taskIndexById := make(map[int64]int)
	newTasks := make([]task.Task, len(tasks.Items))
	for k, t := range tasks.Items {
		taskIndexById[t.Id] = k
		nt, err := task.CloneTaskTx(tx, orgUuid, newWf.Id, newWf.Uuid, t)
		newWf.ElementsLayout = replaceUuid(newWf.ElementsLayout, t.Uuid, nt.Uuid)
		if err != nil {
			log.Error(err)
			return uuid.Nil.String(), err
		}
		newTasks[k] = *nt
	}

	for _, tt := range triggers.Items {
		ntt := tt
		ntt.Id = 0
		ntt.Uuid = uuid.NewString()
		ntt.Schedule = tt.Schedule
		ntt.IsActive = tt.IsActive
		ntt.PipelineId = newWf.Id
		ntt.PipelineUuid = newWf.Uuid
		ntt.TriggerType = tt.TriggerType

		if tt.TriggerTaskId != 0 {
			triggerTaskIndex := taskIndexById[tt.TriggerTaskId]
			ntt.TriggerTaskId = newTasks[triggerTaskIndex].Id
			ntt.TriggerTaskUuid = newTasks[triggerTaskIndex].Uuid
		}

		triggeredTaskIndex := taskIndexById[tt.TriggeredTaskId]
		ntt.TriggeredTaskId = newTasks[triggeredTaskIndex].Id
		ntt.TriggeredTaskUuid = newTasks[triggeredTaskIndex].Uuid

		newWf.ElementsLayout = replaceUuid(newWf.ElementsLayout, tt.Uuid, ntt.Uuid)
		newWf.ElementsLayout = replaceUuid(newWf.ElementsLayout, tt.TriggerTaskUuid, ntt.TriggerTaskUuid)
		newWf.ElementsLayout = replaceUuid(newWf.ElementsLayout, tt.TriggeredTaskUuid, ntt.TriggeredTaskUuid)

		_, _, err := tasktrigger.CreateTaskTriggerTx(tx, &ntt)
		if err != nil {
			return "", err
		}
	}

	err = pipeline.UpdateElementsLayoutTx(tx, &newWf)
	if err != nil {
		return "", err
	}

	return newWf.Uuid, nil
}

func replaceUuid(elementsLayout, search, replace string) string {
	return strings.ReplaceAll(elementsLayout, search, replace)
}

func CreateSharedPipeline(db *sql.DB, tpl *pipeline.Pipeline, creatorUuid string, orgUuid string) (*SharedPipeline, error) {
	sl := &SharedPipeline{}
	sl.PipelineUuid = tpl.Uuid
	sl.OrganizationUuid = orgUuid
	sl.CreatorUuid = creatorUuid
	sl.IsActive = 1
	sl.IsLinkPublished = 1
	sl.ShareLink = shortuuid.New()

	q := `INSERT INTO shared_pipelines(pipeline_uuid, organization_uuid, creator_uuid, share_link, is_active, is_link_published) 
		  VALUES
			(?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(q, sl.PipelineUuid, sl.OrganizationUuid, sl.CreatorUuid, sl.ShareLink, sl.IsActive, sl.IsLinkPublished)
	if err != nil {
		return nil, err
	}

	sl.Id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func IsPipelineShared(db *sql.DB, pipelineUuid string) (bool, error) {
	q := `SELECT COUNT(*) FROM shared_pipelines WHERE pipeline_uuid = ? AND is_active = 1`
	row := db.QueryRow(q, pipelineUuid)
	if row.Err() != nil {
		return false, row.Err()
	}

	val := 0
	err := row.Scan(&val)
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

func FindActiveShareLink(db *sql.DB, shareLink string) (*SharedPipeline, error) {
	q := `SELECT ` + sharedPipelineSqlFields() + ` FROM shared_pipelines WHERE share_link = ? AND is_active = 1 AND is_link_published = 1`
	return sharedPipelineFromRow(db.QueryRow(q, shareLink))
}

func FindActiveShareLinkForPipeline(db *sql.DB, pipelineUuid string) (*SharedPipeline, error) {
	q := `SELECT ` + sharedPipelineSqlFields() + ` FROM shared_pipelines WHERE pipeline_uuid = ? AND is_active = 1 AND is_link_published = 1`
	return sharedPipelineFromRow(db.QueryRow(q, pipelineUuid))
}

func FindAllActiveShareLinksForUser(db *sql.DB, userUuid string) ([]*SharedPipeline, error) {
	result := make([]*SharedPipeline, 0)
	q := `SELECT ` + sharedPipelineSqlFields() + ` FROM shared_pipelines WHERE creator_uuid = ? AND is_active = 1 AND is_link_published = 1`
	rows, err := db.Query(q, userUuid)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		sl := SharedPipeline{}
		err = rows.Scan(&sl.Id, &sl.PipelineUuid, &sl.OrganizationUuid, &sl.CreatorUuid, &sl.ShareLink, &sl.IsActive, &sl.IsLinkPublished, &sl.CreatedAt, &sl.UpdatedAt)
		if err != nil {
			return result, err
		}

		result = append(result, &sl)
	}

	return result, nil
}

func DeactivateSharedPipeline(db *sql.DB, pipelineUuid string) error {
	q := `UPDATE shared_pipelines SET is_active = 0 WHERE pipeline_uuid = ?`
	_, err := db.Exec(q, pipelineUuid)
	if err != nil {
		return err
	}

	return nil
}

func SetShareLinkPublished(db *sql.DB, pipelineUuid string, isLinkPublished bool) error {
	q := `UPDATE shared_pipelines SET is_link_published = ? WHERE pipeline_uuid = ? AND is_active = 1`
	isLinkPublishedParam := 0
	if isLinkPublished {
		isLinkPublishedParam = 1
	}
	_, err := db.Exec(q, isLinkPublishedParam, pipelineUuid)
	if err != nil {
		return err
	}

	return nil
}

func sharedPipelineSqlFields() string {
	return "id, pipeline_uuid, organization_uuid, creator_uuid, share_link, is_active, is_link_published, created_at, updated_at"
}

func sharedPipelineFromRow(row *sql.Row) (*SharedPipeline, error) {
	if row.Err() != nil {
		return nil, row.Err()
	}

	sl := SharedPipeline{}
	err := row.Scan(
		&sl.Id,
		&sl.PipelineUuid,
		&sl.OrganizationUuid,
		&sl.CreatorUuid,
		&sl.ShareLink,
		&sl.IsActive,
		&sl.IsLinkPublished,
		&sl.CreatedAt,
		&sl.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

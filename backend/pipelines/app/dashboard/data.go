package dashboard

import (
	"time"
	"database/sql"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/task"

	log "github.com/sirupsen/logrus"
)

func GetDashboardByOrganizationUuid(db *sql.DB, uuid string, userUuid string) (*Dashboard, error) {
	var dashboard Dashboard

	pipelines, err := pipeline.GetPipelinesByOrganizationUuid(db, uuid)

	if err != nil {
		return nil, err
	}

	if pipelines == nil {
		dashboard.NumActivePipelines = 0
		dashboard.NumActiveMetrics = 0
		dashboard.NumNewPipelines = 0
		dashboard.NumNewMetrics = 0
		dashboard.NumRecentlyUpdatedPipelines = 0
		dashboard.NumRecentlyUpdatedMetrics = 0
		dashboard.NumMyPipelines = 0
		dashboard.NumMyMetrics = 0
		dashboard.NumScheduledPipelines = 0
		dashboard.NumExtTriggeredPipelines = 0
		dashboard.NumPipelineTemplates = 0
		dashboard.NumMetricTemplates = 0
	} else {
		var NumActivePipelines int
		var NumActiveMetrics int
		var NumNewPipelines int
		var NumNewMetrics int
		var NumRecentlyUpdatedPipelines int
		var NumRecentlyUpdatedMetrics int
		var NumScheduledPipelines int
		var NumMyPipelines int
		var NumMyMetrics int
		var NumPipelineTemplates int
		var NumMetricTemplates int
		var createdAt time.Time
		var updatedAt time.Time

		t := time.Now().Add(time.Hour * 24 * -30)

		for i := 0; i < len(pipelines.Items); i++ {
		    if pipelines.Items[i].Type == "generic" {
				NumActivePipelines++

				createdAt, err = time.Parse(time.RFC3339, pipelines.Items[i].CreatedAt)
				if err == nil {
					if t.Before(createdAt) {
						NumNewPipelines++
					}
				}

				updatedAt, err = time.Parse("2006-01-02 15:04:05", pipelines.Items[i].UpdatedAt)
				if err == nil {
					if t.Before(updatedAt) {
						NumRecentlyUpdatedPipelines++
					}
				}

				if pipelines.Items[i].Schedule != "" {
					NumScheduledPipelines++
				}

				if pipelines.Items[i].IsTemplate == 1 {
					NumPipelineTemplates++
				}

				if pipelines.Items[i].CreatorUuid == userUuid {
					NumMyPipelines++
				}
			} else if pipelines.Items[i].Type == "metric" {
				NumActiveMetrics++

				createdAt, err = time.Parse(time.RFC3339, pipelines.Items[i].CreatedAt)
				if err == nil {
					if t.Before(createdAt) {
						NumNewMetrics++
					}
				}

				updatedAt, err = time.Parse("2006-01-02 15:04:05", pipelines.Items[i].UpdatedAt)
				if err == nil {
					if t.Before(updatedAt) {
						NumRecentlyUpdatedMetrics++
					}
				}

				if pipelines.Items[i].IsTemplate == 1 {
					NumMetricTemplates++
				}

				if pipelines.Items[i].CreatorUuid == userUuid {
					NumMyMetrics++
				}
			}
		}

		dashboard.NumActivePipelines = NumActivePipelines
		dashboard.NumActiveMetrics = NumActiveMetrics
		dashboard.NumNewPipelines = NumNewPipelines
		dashboard.NumNewMetrics = NumNewMetrics
		dashboard.NumRecentlyUpdatedPipelines = NumRecentlyUpdatedPipelines
		dashboard.NumRecentlyUpdatedMetrics = NumRecentlyUpdatedMetrics
		dashboard.NumScheduledPipelines = NumScheduledPipelines
		dashboard.NumMyPipelines = NumMyPipelines
		dashboard.NumMyMetrics = NumMyMetrics
		dashboard.NumPipelineTemplates = NumPipelineTemplates
		dashboard.NumMetricTemplates = NumMetricTemplates

		eW, _ := task.GetExternallyTriggeredPipelineCount(db, uuid)
		dashboard.NumExtTriggeredPipelines = eW
	}

	return &dashboard, nil
}

func GetGroupedItemsByOrganizationUuid(db *sql.DB, uuid string, itemType string, groupBy string, userUuid string) (*GroupedItems, error) {
	var Query string

	if groupBy == GroupByMonth {
		Query = `SELECT
			COUNT(id),
			YEAR(created_at), 
			DATE_FORMAT(created_at, '%b'),
			0
		FROM pipelines
		WHERE
			organization_uuid = ?
			AND is_active = 1
			AND type = ?
		GROUP BY YEAR(created_at), MONTH(created_at)
		ORDER BY YEAR(created_at), MONTH(created_at)`
    } else {
    	Query = `SELECT
			COUNT(id),
			YEAR(created_at), 
			DATE_FORMAT(created_at, '%b'),
			WEEK(created_at)
		FROM pipelines
		WHERE
			organization_uuid = ?
			AND is_active = 1
			AND type = ?
		GROUP BY YEAR(created_at), MONTH(created_at), WEEK(created_at)
		ORDER BY YEAR(created_at), MONTH(created_at), WEEK(created_at)`
	}

	stmt, err := db.Prepare(Query)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid, itemType)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var items GroupedItems

	for rows.Next() {
		var item GroupedItem
		err := rows.Scan(&item.Count, &item.Year, &item.Month, &item.Week)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		items.Items = append(items.Items, item)
	}

	if err := rows.Err(); err != nil {
		log.Error(err)
		return &items, err
	}

	return &items, err
}

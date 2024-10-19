package migration

import (
	"ylem_statistics/config"
	"ylem_statistics/services/db"
)

func init() {
	db.AddMigration(
		20231113203511,

		`CREATE TABLE `+config.Cfg().DB.StatsTable+`(
			uuid UUID,
			executor_uuid UUID,
			organization_uuid UUID,
			creator_uuid UUID,
			pipeline_uuid UUID,
			pipeline_run_uuid UUID,
			task_uuid UUID,
			task_type String,
			output BLOB DEFAULT '',
			pipeline_type String DEFAULT '',
			metric_value Float64 DEFAULT 0,
			is_metric_value_set UInt8  DEFAULT 0,
			is_initial_task UInt8,
			is_final_task UInt8,
			is_successful UInt8,
			is_fatal_failure UInt8,
			executed_at DateTime64(6, 'UTC'),
			duration UInt32,
			PRIMARY KEY(uuid)
		) ENGINE = MergeTree
		ORDER BY (uuid, executed_at)
		`,

		`DROP TABLE `+config.Cfg().DB.StatsTable,
	)
}

package provider

import (
	"database/sql"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/helpers"
)

type PipelineProvider interface {
	GetPipeline(id int64) (*pipeline.Pipeline, error)
}

type DbPipelineProvider struct {
	db *sql.DB
}

func (wp *DbPipelineProvider) GetPipeline(id int64) (*pipeline.Pipeline, error) {
	return pipeline.GetPipelineById(wp.db, id)
}

func NewPipelineProvider() *DbPipelineProvider {
	return &DbPipelineProvider{
		db: helpers.DbConn(),
	}
}

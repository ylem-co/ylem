package provider

import (
	"context"
	"sync"
	"time"
	"database/sql"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/pipeline/run"

	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
)

type TaskProvider interface {
	GetTask(id int64, config run.PipelineRunConfig) (*task.Task, bool, error)
	GetInitialTasks(pipelineUuid string, config run.PipelineRunConfig) ([]*task.Task, error)
	IsFinalTask(uid string, config run.PipelineRunConfig) (bool, error)
}

type DbTaskProvider struct {
	Db *sql.DB
}

func (p *DbTaskProvider) GetTask(id int64, config run.PipelineRunConfig) (*task.Task, bool, error) {
	t, err := task.GetTaskById(p.Db, id)
	if err != nil {
		return nil, false, err
	}
	isFinalTask, err := p.IsFinalTask(t.Uuid, config)

	return t, isFinalTask, err // task, is final, error
}

func (p *DbTaskProvider) GetInitialTasks(pipelineUuid string, config run.PipelineRunConfig) ([]*task.Task, error) {
	uid, err := uuid.Parse(pipelineUuid)
	if err != nil {
		return nil, err
	}
	tasks, err := task.GetInitialTasks(p.Db, uid, config)
	if err != nil {
		return nil, err
	}

	tasksToTrigger := make([]*task.Task, 0)
	tasksToTrigger = append(tasksToTrigger, tasks...)

	return tasksToTrigger, nil
}

func (p *DbTaskProvider) IsFinalTask(uid string, config run.PipelineRunConfig) (bool, error) {
	tUid, err := uuid.Parse(uid)
	if err != nil {
		return false, err
	}
	ids, err := tasktrigger.GetTriggeredTaskIds(p.Db, tUid, "", config)

	return len(ids) == 0, err
}

type CacheItem struct {
	task    *task.Task
	isFinal bool
}

type CachingTaskProvider struct {
	InnerProvider TaskProvider
	Ctx           context.Context
	taskCache     *lru.Cache
	mu            *sync.RWMutex
}

func (p *CachingTaskProvider) Init() error {
	var err error
	p.taskCache, err = lru.New(10000)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				p.taskCache, err = lru.New(10000)
				if err != nil {
					return
				}

			case <-p.Ctx.Done():
				return
			}
		}
	}()

	p.mu = &sync.RWMutex{}

	return nil
}

func (p *CachingTaskProvider) GetTask(id int64, config run.PipelineRunConfig) (*task.Task, bool, error) {
	if ci, ok := p.taskCache.Get(id); ok {
		return ci.(CacheItem).task, ci.(CacheItem).isFinal, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if ci, ok := p.taskCache.Get(id); ok {
		return ci.(CacheItem).task, ci.(CacheItem).isFinal, nil
	}

	t, isFinal, err := p.InnerProvider.GetTask(id, config)
	if err != nil {
		return nil, isFinal, err
	}

	newCi := CacheItem{
		task:    t,
		isFinal: isFinal,
	}
	p.taskCache.Add(id, newCi)

	return newCi.task, newCi.isFinal, err
}

func (p *CachingTaskProvider) GetInitialTasks(pipelineUuid string, config run.PipelineRunConfig) ([]*task.Task, error) {
	return p.InnerProvider.GetInitialTasks(pipelineUuid, config)
}
func (p *CachingTaskProvider) IsFinalTask(uid string, config run.PipelineRunConfig) (bool, error) {
	return p.InnerProvider.IsFinalTask(uid, config)
}

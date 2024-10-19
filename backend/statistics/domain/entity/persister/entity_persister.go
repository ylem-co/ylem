package persister

import (
	"sync"
	"ylem_statistics/domain/entity"
	"ylem_statistics/services/db"

	"gorm.io/gorm"
)

type EntityPersister interface {
	CreateTaskRun(tr *entity.TaskRun) error
}

type entityPersister struct {
	db *gorm.DB
}

func (p *entityPersister) CreateTaskRun(tr *entity.TaskRun) error {
	result := p.db.Create(tr)

	return result.Error
}

var instance *entityPersister
var mu = &sync.RWMutex{}

func Instance() EntityPersister {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		instance = newEntityPersister()
	}

	return instance
}

func newEntityPersister() *entityPersister {
	dbInstance, err := db.Instance()
	if err != nil {
		panic(err)
	}
	return &entityPersister{
		db: dbInstance,
	}
}

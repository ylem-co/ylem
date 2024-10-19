package messaging

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
	"ylem_pipelines/app/task"
	"ylem_pipelines/services/ylem_integrations"

	"github.com/google/uuid"
	messaging "github.com/ylem-co/shared-messaging"
	lru "github.com/hashicorp/golang-lru"
	log "github.com/sirupsen/logrus"
)

type QueryTaskMessageFactory struct {
	ctx         context.Context
	sourceCache *lru.Cache
	mu          *sync.RWMutex
	client      ylem_integrations.Client
}

func (f *QueryTaskMessageFactory) init() error {
	f.mu = &sync.RWMutex{}
	err := f.initSourceCache()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			// reset source cache every 30 seconds
			case <-time.After(time.Second * 30):
				err := f.initSourceCache()
				if err != nil {
					panic(err)
				}

			case <-f.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (f *QueryTaskMessageFactory) initSourceCache() error {
	var err error
	f.sourceCache, err = lru.New(10000)
	log.Trace("Source cache reset")

	return err
}

func (f *QueryTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.Query)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.Query{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	msg := &messaging.RunQueryTask{
		Task:  task,
		Query: impl.SQLQuery,
	}

	sUid, err := uuid.Parse(impl.SourceUuid)
	if err != nil {
		return nil, err
	}
	s, err := f.getSQLIntegration(sUid)
	if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
		return nil, NewErrorRepeatable("integration service is unavailable")
	} else if err != nil {
		return nil, err
	}
	msg.Source = *s

	return messaging.NewEnvelope(msg), nil
}

func (f *QueryTaskMessageFactory) getSQLIntegration(uid uuid.UUID) (*messaging.SQLIntegration, error) {
	uidStr := uid.String()
	if s, ok := f.sourceCache.Get(uidStr); ok {
		return s.(*messaging.SQLIntegration), nil
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if s, ok := f.sourceCache.Get(uidStr); ok {
		return s.(*messaging.SQLIntegration), nil
	}

	s, err := f.client.GetSQLIntegration(uid)
	if err != nil {
		return s, err
	}

	f.sourceCache.Add(uidStr, s)

	return s, nil
}

func NewRunQueryTaskMessageFactory(ctx context.Context) (*QueryTaskMessageFactory, error) {
	ycl, err := ylem_integrations.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	f := &QueryTaskMessageFactory{
		ctx:    ctx,
		client: ycl,
	}
	err = f.init()
	if err != nil {
		return nil, err
	}

	return f, nil
}

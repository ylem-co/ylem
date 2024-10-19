package schedule

import (
	"context"
	"strconv"
	"sync"
	"time"
	"database/sql"
	"encoding/json"
	"ylem_pipelines/app/envvariable"
	"ylem_pipelines/app/schedule"
	"ylem_pipelines/config"
	"ylem_pipelines/helpers"
	msgsrv "ylem_pipelines/services/messaging"
	"ylem_pipelines/services/provider"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

const PrePublishPeriod = time.Duration(0)

type PipelineRunContextProvider interface {
	GetPipelineRunContexts() (<-chan msgsrv.PipelineRunContext, error)
}

type ScheduledPipelineRunContextProvider struct {
	db               *sql.DB
	ctx              context.Context
	pipelineProvider provider.PipelineProvider
	taskProvider     provider.TaskProvider
	prePublishPeriod time.Duration
}

func (p *ScheduledPipelineRunContextProvider) GetPipelineRunContexts() (<-chan msgsrv.PipelineRunContext, error) {
	result := make(chan msgsrv.PipelineRunContext)
	srs, tx, err := schedule.GetSchedulesForPublishing(p.db, p.prePublishPeriod)
	if err != nil || len(srs) == 0 {
		close(result)

		if tx != nil {
			_ = tx.Rollback()
		}

		return result, err
	}

	log.Tracef("Loaded %d scheduled runs", len(srs))

	go func() {
		defer close(result)
		defer tx.Commit()  //nolint:all

		for _, sr := range srs {
			select {
			case <-p.ctx.Done():
				// if stop signal received, close channel and stop
				return

			default:
				wrUid := uuid.New()
				if sr.PipelineRunUuid != uuid.Nil {
					wrUid = sr.PipelineRunUuid
				}

				wf, err := p.pipelineProvider.GetPipeline(sr.PipelineId)
				if err != nil {
					log.Error(err)
					continue
				}

				if wf == nil {
					log.Debugf("Pipeline not found, deleting scheduled run %d", sr.Id)
					_, _ = schedule.DeleteTx(tx, sr.Id)
					continue
				}

				sysVars, err := envvariable.GetEnvVariablesByOrganizationUuidTx(tx, wf.OrganizationUuid)
				if err != nil {
					log.Error(err)
					continue
				}

				if sr.EnvVars == nil {
					sr.EnvVars = make(map[string]interface{})
				}

				for _, sv := range sysVars.Items {
					if _, ok := sr.EnvVars[sv.Name]; ok {
						continue
					}

					sr.EnvVars[sv.Name] = sv.Value
				}

				tasks, err := p.taskProvider.GetInitialTasks(wf.Uuid, sr.Config)
				if err != nil {
					continue
				}

				wrc := msgsrv.PipelineRunContext{
					ScheduledRunId:  sr.Id,
					PipelineRunUuid: wrUid,
					PipelineUuid:    wf.Uuid,
					PipelineType:    wf.Type,
				}

				for _, t := range tasks {
					isFinal, err := p.taskProvider.IsFinalTask(t.Uuid, sr.Config)
					if err != nil {
						continue
					}

					wrc.OrganizationUuid = t.OrganizationUuid
					input := sr.Input
					if len(input) == 0 {
						input, _ = json.Marshal(nil)
					}

					wrc.TaskRuns = append(wrc.TaskRuns, msgsrv.TaskRunContext{
						PipelineRunContext: wrc,
						Task:               t,
						IsInitialTask:      true,
						IsFinalTask:        isFinal,
						ExecuteAt:          sr.ExecuteAt,
						PipelineRunUuid:    wrUid,
						Input:              input,
						Meta: messaging.Meta{
							EnvVars: sr.EnvVars,
							PipelineRunConfig: messaging.PipelineRunConfig{
								TaskIds: messaging.IdList{
									Type: sr.Config.TaskIds.Type,
									Ids:  sr.Config.TaskIds.Ids,
								},
								TaskTriggerIds: messaging.IdList{
									Type: sr.Config.TaskTriggerIds.Type,
									Ids:  sr.Config.TaskTriggerIds.Ids,
								},
							},
						},
					})
				}

				result <- wrc
			}
		}
	}()

	return result, nil
}

type Publisher struct {
	topic                      string
	brokers                    []string
	pipelineRunContextProvider PipelineRunContextProvider
	messageFactory             msgsrv.MessageFactory
	db                         *sql.DB
	ctx                        context.Context
	emitter                    *goka.Emitter
	stopped                    bool
	mu                         *sync.Mutex
}

func (p *Publisher) Start() error {
	go func() {
		// stop publisher when Done() is closed
		<-p.ctx.Done()
		log.Info("Stopping schedule publisher")
		p.stopped = true
	}()

	return p.runCycle()
}

func (p *Publisher) initEmitter() error {
	var err error
	p.emitter, err = goka.NewEmitter(p.brokers, goka.Stream(p.topic), new(messaging.MessageCodec))

	return err
}

func (p *Publisher) runCycle() error {
	log.Info("Schedule publisher started")

	for {
		publishedNum, err := p.run()
		log.Debugf("Published %d messages", publishedNum)
		if err != nil {
			return err
		}

		if publishedNum == 0 {
			log.Debug("Nothing to do, waiting...")
			select {
			case <-time.After(time.Second):
				continue

			case <-p.ctx.Done():
				return nil
			}
		}

		if p.stopped {
			return nil
		}
	}
}

func (p *Publisher) run() (int, error) {
	err := p.initEmitter()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	wrcsChan, err := p.pipelineRunContextProvider.GetPipelineRunContexts()
	if err != nil {
		log.Error(err)

		return 0, err
	}

	wrcs := make([]msgsrv.PipelineRunContext, 0)
	for wrc := range wrcsChan {
		wrcs = append(wrcs, wrc)
	}

	log.Trace("Starting publishing")
	tx, err := p.db.Begin()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	publishedNum, err := p.PublishSchedules(tx, wrcs)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return publishedNum, err
	}

	err = p.emitter.Finish()
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return publishedNum, err
	}

	log.Trace("Committing transaction")
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return publishedNum, err
	}

	return publishedNum, nil
}

func (p *Publisher) PublishSchedules(tx *sql.Tx, wrcs []msgsrv.PipelineRunContext) (int, error) {
	i := 0

	for _, wrc := range wrcs {
		orgUid, err := uuid.Parse(wrc.OrganizationUuid)
		if err != nil {
			log.Error(err)
			_, _ = schedule.Delete(p.db, wrc.ScheduledRunId)
			return i, err
		}

		for _, trc := range wrc.TaskRuns {
			if p.stopped {
				return i, nil
			}

			published, err := p.publishOne(tx, wrc, trc)
			if err != nil {
				return i, err
			}

			if published {
				i++
			}
		}

		err = schedule.IncrementCurrentPipelineRunCount(tx, orgUid, wrc.PipelineType)
		if err != nil {
			log.Error(err)
		}
	}

	return i, nil
}

func (p *Publisher) publishOne(tx *sql.Tx, wrc msgsrv.PipelineRunContext, trc msgsrv.TaskRunContext) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovering from panic in Publisher.publishOne(): %v", r)
			_, _ = schedule.Delete(p.db, wrc.ScheduledRunId)
		}
	}()
	var err error

	msg, err := p.messageFactory.CreateMessage(trc)
	if err != nil {
		if _, ok := err.(msgsrv.ErrorRepeatable); ok {
			log.Errorf("repeatable error, scheduled run will be retried: %s", err)
			return false, nil
		}

		log.Error(err)
		_, _ = schedule.Delete(p.db, wrc.ScheduledRunId)
		return false, nil
	}

	prom, err := p.emitter.EmitWithHeaders(
		trc.PipelineRunUuid.String(),
		msg,
		GetMessageHeaders(p.topic, trc),
	)
	if err != nil {
		return false, err
	}

	prom.Then(func(err error) {
		p.mu.Lock()
		defer p.mu.Unlock()
		if err != nil {
			log.Errorf("message publishing failed: %s", err)
			return
		}

		log.Tracef("Message published, deleting scheduled run %d", wrc.ScheduledRunId)
		_, _ = schedule.Delete(p.db, wrc.ScheduledRunId)
	})

	return true, nil
}

func GetMessageHeaders(topic string, trc msgsrv.TaskRunContext) goka.Headers {
	h := make(goka.Headers)
	h["scheduler-epoch"] = []byte(strconv.FormatInt(trc.ExecuteAt.Unix(), 10))
	h["scheduler-target-topic"] = []byte(topic)
	h["scheduler-target-key"] = []byte(trc.Task.Uuid)

	return h
}

func NewPublisher(ctx context.Context) (*Publisher, error) {
	db := helpers.DbConn()
	tp := &provider.DbTaskProvider{
		Db: db,
	}

	ctp := &provider.CachingTaskProvider{
		InnerProvider: tp,
		Ctx:           ctx,
	}

	err := ctp.Init()
	if err != nil {
		return nil, err
	}

	wrcp := &ScheduledPipelineRunContextProvider{
		db:               db,
		pipelineProvider: provider.NewPipelineProvider(),
		taskProvider:     ctp,
		prePublishPeriod: PrePublishPeriod,
		ctx:              ctx,
	}

	return NewPublisherWithProvider(ctx, wrcp)
}

func NewPublisherWithProvider(ctx context.Context, wrcp PipelineRunContextProvider) (*Publisher, error) {
	cfg := config.Cfg().Kafka
	db := helpers.DbConn()

	mf, err := msgsrv.NewCompositeMessageFactory(ctx)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		ctx:                        ctx,
		brokers:                    cfg.BootstrapServers,
		topic:                      cfg.TaskRunsTopic,
		pipelineRunContextProvider: wrcp,
		messageFactory:             mf,
		db:                         db,
		mu:                         &sync.Mutex{},
	}, nil
}

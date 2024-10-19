package command

import (
	"time"
	"ylem_pipelines/app/schedule"
	"ylem_pipelines/helpers"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

const TimeQuantum = time.Minute
const PregenerationPeriod = time.Hour * 48

var scheduleGeneratorStartHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Schedule generator started")
	db := helpers.DbConn()
	defer db.Close()

	for {
		srs, tx, err := schedule.GetLastScheduledRunsForScheduling(db, PregenerationPeriod)
		if err != nil {
			return err
		}
		if len(srs) == 0 {
			log.Debug("Nothing to do, waiting...")
			err = tx.Commit()
			if err != nil {
				return err
			}
			select {
			case <-time.After(time.Second * 10):
				continue

			case <-c.Done():
				return nil
			}
		}

		if err != nil {
			_ = tx.Rollback()
			return err
		}

		for _, sr := range srs {
			log.Tracef("Generating schedule for pipeline %d", sr.PipelineId)
			var start time.Time
			if sr.ExecuteAt == nil {
				log.Tracef("No schedules found for pipeline %d, generating from now.", sr.PipelineId)
				start = time.Now().Add(TimeQuantum)
			} else {
				start = sr.ExecuteAt.Add(TimeQuantum)
				log.Tracef("Generating schedues for pipeline %d from %s.", sr.PipelineId, start.Format(helpers.DB_TIME_TIMESTAMP))
			}

			end := time.Now().Add(PregenerationPeriod)
			if start.After(end) {
				log.Debug("Start after end")
				continue
			}

			newSrs := generateSchedule(sr.PipelineId, sr.Schedule, start, end)
			err = schedule.AddScheduledRunsTx(tx, newSrs)
			if err != nil {
				return err
			}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}
}

func generateSchedule(pipelineId int64, expr string, start time.Time, end time.Time) []schedule.ScheduledRun {
	result := make([]schedule.ScheduledRun, 0)
	currentTime := start
	gron := gronx.New()
	for currentTime.Before(end) || currentTime.Equal(end) {
		cTime := currentTime.Truncate(TimeQuantum)
		isDue, err := gron.IsDue(expr, cTime)
		if err != nil {
			return result
		}
		if isDue {
			result = append(result, schedule.ScheduledRun{
				PipelineId:      pipelineId,
				PipelineRunUuid: uuid.New(),
				ExecuteAt:       &cTime,
			})
		}
		currentTime = currentTime.Add(TimeQuantum)
	}

	return result
}

var ScheduleGeneratorStart = &cli.Command{
	Name:   "start",
	Usage:  "Start schedule generator",
	Action: scheduleGeneratorStartHandler,
}

var ScheduleGeneratorCommands = &cli.Command{
	Name:  "schedulegen",
	Usage: "Schedule generator commands",
	Subcommands: []*cli.Command{
		ScheduleGeneratorStart,
	},
}

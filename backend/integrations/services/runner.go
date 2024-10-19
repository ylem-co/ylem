package services

import (
	log "github.com/sirupsen/logrus"
	"reflect"

	messaging "github.com/ylem-co/shared-messaging"
)

func ProcessMessage(task interface{}) {
	switch task := task.(type) {
	case *messaging.TaskRunResult:
		// do something
	default:
		log.Debugf(`Ignoring the message of type "%s"`, reflect.TypeOf(task).String())
	}
}

/*func runMeasured(f func() *messaging.TaskRunResult) *messaging.TaskRunResult {
	start := time.Now()
	tr := f()
	tr.ExecutedAt = time.Now()
	tr.Duration = time.Since(start) / time.Millisecond

	return tr
}*/

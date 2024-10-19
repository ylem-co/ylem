package kafka

import (
	"encoding/json"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	messaging "github.com/ylem-co/shared-messaging"
)

func DecodeKafkaTaskValue(t messaging.Task, messageName string, tr *messaging.TaskRunResult) (interface{}, error) {
	var (
		err       error
		taskValue interface{}
	)

	if len(t.Input) == 0 {
		return make(map[string]interface{}), nil
	}

	err = json.Unmarshal(t.Input, &taskValue)

	if err != nil {
		HandleBadRequestError(t.TaskUuid, messageName, err, tr)

		return taskValue, err
	}

	return taskValue, nil
}

func HandleBadRequestError(taskUuid uuid.UUID, messageName string, err error, tr *messaging.TaskRunResult) {
	log.Errorf(
		`could not execute task "%s"" with uuid "%s": %v`,
		messageName,
		taskUuid,
		err,
	)

	tr.IsSuccessful = false
	tr.Errors = []messaging.TaskRunError{
		{
			Code:     messaging.ErrorBadRequest,
			Severity: messaging.ErrorSeverityError,
			Message:  err.Error(),
		},
	}
}

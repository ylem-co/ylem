package messaging

import (
	"context"
	"fmt"
	"reflect"
	"ylem_pipelines/app/task"
	"ylem_pipelines/helpers"
	"ylem_pipelines/services/ylem_integrations"

	"github.com/google/uuid"
	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
)

type ApiTaskMessageFactory struct {
	ctx         context.Context
	integrationsClient ylem_integrations.Client
}

func (f *ApiTaskMessageFactory) CreateMessage(trc TaskRunContext) (*messaging.Envelope, error) {
	t := trc.Task
	impl, ok := t.Implementation.(*task.ApiCall)
	if !ok {
		return nil, fmt.Errorf(
			"wrong task type. Expected %s, got %s",
			reflect.TypeOf(&task.ApiCall{}).String(),
			reflect.TypeOf(t.Implementation).String(),
		)
	}

	dUid, err := uuid.Parse(impl.DestinationUuid)
	if err != nil {
		return nil, err
	}

	d, err := f.integrationsClient.GetApiIntegration(dUid)
	if _, ok := err.(ylem_integrations.ErrorServiceUnavailable); ok {
		return nil, NewErrorRepeatable("Ylem_integrations service is unavailable")
	} else if err != nil {
		return nil, err
	}

	headers, err := helpers.DecodeHeaders(impl.Headers)
	if err != nil {
		log.Error(err)
	}
	task, err := createTaskMessage(trc)
	if err != nil {
		return nil, err
	}

	dm, err := createIntegrationMessage(d.Integration)
	if err != nil {
		return nil, err
	}
	msg := &messaging.CallApiTask{
		Task: task,
		Integration: messaging.ApiIntegration{
			Integration:           dm,
			Method:                d.Method,
			AuthType:              d.AuthType,
			AuthBearerToken:       d.AuthBearerToken,
			AuthBasicUserName:     d.AuthBasicUserName,
			AuthBasicUserPassword: d.AuthBasicUserPassword,
			AuthHeaderName:        d.AuthHeaderName,
			AuthHeaderValue:       d.AuthHeaderValue,
		},
		Type:             impl.Type,
		Payload:          impl.Payload,
		QueryString:      impl.QueryString,
		Headers:          headers,
		Severity:         t.Severity,
		AttachedFileName: impl.AttachedFileName,
	}

	return messaging.NewEnvelope(msg), nil
}

func NewApiTaskMessageFactory(ctx context.Context) (*ApiTaskMessageFactory, error) {
	ycl, err := ylem_integrations.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	f := &ApiTaskMessageFactory{
		ctx:         ctx,
		integrationsClient: ycl,
	}

	return f, nil
}

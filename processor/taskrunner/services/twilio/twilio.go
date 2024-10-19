package twilio

import (
	"errors"
	"ylem_taskrunner/config"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(ToPhoneNumber string, Text string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Cfg().Twilio.AccountSid,
		Password: config.Cfg().Twilio.AuthToken,
	})

	if client == nil {
		return errors.New("could not create Twilio client")
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(ToPhoneNumber)
	params.SetFrom(config.Cfg().Twilio.NumberFrom)
	params.SetBody(Text)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}

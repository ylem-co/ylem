package twilio

import (
	"fmt"
	"errors"
	"encoding/json"
	"ylem_taskrunner/config"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsAppPayload struct {
    PhoneTo          string `json:"phoneTo"`
    ContentVariables string `json:"contentVariables"`
}

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

func SendWhatsAppMessage(ContentSid string, FromPhoneNumber string, Text string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Cfg().Twilio.AccountSid,
		Password: config.Cfg().Twilio.AuthToken,
	})

	if client == nil {
		return errors.New("could not create Twilio client")
	}

	data := WhatsAppPayload{}
    err := json.Unmarshal([]byte(Text), &data)
    if err != nil {
		return err
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(fmt.Sprintf("%s%s", "whatsapp:", FromPhoneNumber))
	params.SetTo(fmt.Sprintf("%s%s", "whatsapp:", data.PhoneTo))
	params.SetContentSid(ContentSid)
	params.SetContentVariables(data.ContentVariables)

	_, err = client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}

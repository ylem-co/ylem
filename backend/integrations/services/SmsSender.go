package services

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"ylem_integrations/config"
)

func SendSms(ToPhoneNumber string, Text string) error {
	var config config.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Twilio.AccountSid,
		Password: config.Twilio.AuthToken,
	})

	if client == nil {
		log.Println("Could not create Twilio client")

		return errors.New("could not create Twilio client")
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(ToPhoneNumber)
	params.SetFrom(config.Twilio.NumberFrom)
	params.SetBody(Text)

	_, err = client.Api.CreateMessage(params)
	if err != nil {
		log.Println(err.Error())

		return err
	}

	return nil
}

func SendPhoneNumberVerificationSms(ToPhoneNumber string, Code string) error {
	return SendSms(ToPhoneNumber, "Your Ylem verification code " + Code)
}

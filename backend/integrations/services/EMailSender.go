package services

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/kelseyhightower/envconfig"
	"html/template"
	"log"
	"strings"
	"ylem_integrations/config"
	"ylem_integrations/resources"
)

const (
	EmailDefaultCharset = "UTF-8"
	EmailDefaultSender = "Ylem <admin@ylem.co>"
)

func SendEmail(Recipient string, Subject string, HtmlBody string) (*ses.SendEmailOutput, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region:aws.String("eu-central-1"),
		},
	)
	if err != nil {
		return nil, err
	}

	// Create an SES session.
	service := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
			},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(EmailDefaultCharset),
					Data:    aws.String(HtmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(EmailDefaultCharset),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(EmailDefaultSender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	result, err := service.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr
		}

		return nil, err
	}

	return result, nil
}

func SendEmailConfirmationEmail(Recipient string, DestinationUuid string, ConfirmationCode string) (*ses.SendEmailOutput, error) {
	page := resources.EmbeddedHtmlTemplates["confirmation_email"]
	tpl,_ := template.ParseFS(resources.EmbeddedFileSystem, page)

	var config config.Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	var buffer bytes.Buffer

	link := strings.Replace(config.NetworkConfig.EmailConfirmationUrl, "{uuid}", DestinationUuid, 1)
	data := map[string]interface{}{
		"link": link,
		"code": ConfirmationCode,
	}

	if err := tpl.Execute(&buffer, data); err != nil {
		return nil, err
	}

	return SendEmail(Recipient, "Please Confirm Your Ylem Destination Email", buffer.String())
}

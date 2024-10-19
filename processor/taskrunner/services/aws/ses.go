package aws

import (
	"bytes"
	"fmt"
	"io"
	"ylem_taskrunner/config"

	"gopkg.in/gomail.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	EmailDefaultCharset = "UTF-8"
	EmailDefaultSender = "Ylem <admin@ylem.co>"
)

type EmailFile struct {
	Content []byte
	Name    string
}

func SendEmail(Recipient string, Subject string, HtmlBody string, File *EmailFile) (*ses.SendRawEmailOutput, error) {
	cfg := config.Cfg()

	sess, err := session.NewSession(
		&aws.Config{
			Region:aws.String(cfg.Aws.Region),
		},
	)
	if err != nil {
		return nil, err
	}

	service := ses.New(sess)

	msg := gomail.NewMessage()
	msg.SetHeader("From", EmailDefaultSender)
	msg.SetHeader("To", Recipient)
	msg.SetHeader("Subject", Subject)
	msg.SetBody("text/plain", HtmlBody)

	if File != nil {
		msg.Attach(
			fmt.Sprint(File.Name),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(File.Content)

				return err
			}),
		)
	}

	var emailRaw bytes.Buffer
	_, err = msg.WriteTo(&emailRaw)
	if err != nil {
		return nil, err
	}

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: emailRaw.Bytes(),
		},
		Source: aws.String(EmailDefaultSender),
	}

	result, err := service.SendRawEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr
		}

		return nil, err
	}

	return result, nil
}

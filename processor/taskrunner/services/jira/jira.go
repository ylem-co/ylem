package jira

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/aws/kms"

	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
)

type Issue struct {
	ProjectKey  string
	IssueType   string
	Summary     string
	Description string
}

type Authentication struct {
	CloudId              string
	EncryptedDataKey     []byte
	EncryptedAccessToken []byte
}

func (a *Authentication) Decrypt(ctx context.Context) ([]byte, error) {
	decryptedDataKey, err := kms.DecryptDataKey(
		ctx,
		config.Cfg().Aws.KmsKeyId,
		a.EncryptedDataKey,
	)

	if err != nil {
		return nil, err
	}

	decryptedAccessKey, err := kms.Decrypt(a.EncryptedAccessToken, decryptedDataKey)

	if err != nil {
		return nil, err
	}

	return decryptedAccessKey, nil
}

func CreateTask(ctx context.Context, issue Issue, auth Authentication) error {
	log.Tracef("jira: create task")
	token, err := auth.Decrypt(ctx)
	if err != nil {
		return err
	}

	tp := jira.BearerAuthTransport{
		Token: string(token),
	}
	log.Debugf("jira: decrypted token %s", string(token))
	log.Debugf("jira: auth %v", auth)

	jiraClient, err := jira.NewClient(
		tp.Client(),
		fmt.Sprintf("https://api.atlassian.com/ex/jira/%s/", auth.CloudId),
	)
	if err != nil {
		return err
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Summary: issue.Summary,
			Description: issue.Description,
			Type: jira.IssueType{
				Name: issue.IssueType,
			},
			Project: jira.Project{
				Key: issue.ProjectKey,
			},
		},
	}

	_, resp, err := jiraClient.Issue.CreateWithContext(ctx, &i)
	if err != nil {
		body := []byte("nil")
		if resp != nil {
			body, _ = io.ReadAll(resp.Body)
		}

		return fmt.Errorf("jira issue create: %s; body %s", err.Error(), string(body))
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("jira issue create: %s; body %s", resp.Status, string(body))
	}

	return nil
}

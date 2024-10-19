package opsgenie

import (
	"context"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/aws/kms"

	"github.com/ylem-co/opsgenie-client"
)

type Alert struct {
	Message     string
	Description string
	Priority    string
}

func DecryptKeyAndCreateAlert(ctx context.Context, dataKey []byte, apiKey []byte, alert Alert) error {
	decryptedDataKey, err := kms.DecryptDataKey(
		ctx,
		config.Cfg().Aws.KmsKeyId,
		dataKey,
	)

	if err != nil {
		return err
	}

	decryptedApiKey, err := kms.Decrypt(apiKey, decryptedDataKey)

	if err != nil {
		return err
	}

	opsgenie, _ := opsgenieclient.CreateInstance(ctx, string(decryptedApiKey))

	err = opsgenie.CreateAlert(opsgenieclient.CreateAlertRequest{
		Message:     alert.Message,
		Description: alert.Description,
		Priority:    alert.Priority,
	})

	if err != nil {
		return err
	}

	return nil
}

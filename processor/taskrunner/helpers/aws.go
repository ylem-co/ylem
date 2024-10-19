package helpers

import (
	"context"
	"ylem_taskrunner/config"
	"ylem_taskrunner/services/aws/kms"
)

func DecryptData(ctx context.Context, dataKey []byte, encryptedData []byte) (string, error) {
	decryptedDataKey, err := kms.DecryptDataKey(
		ctx,
		config.Cfg().Aws.KmsKeyId,
		dataKey,
	)

	if err != nil {
		return "", err
	}

	decryptedData, err := kms.Decrypt(encryptedData, decryptedDataKey)

	if err != nil {
		return "", err
	}

	return string(decryptedData), err
}

package kms

import (
	"context"
	"ylem_users/config"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	awstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	log "github.com/sirupsen/logrus"
)

var svc *awskms.Client

func IssueDataKeyWithContext(ctx context.Context) ([]byte, error) {
	masterKeyId := config.Cfg().Aws.KmsKeyId

	if len(masterKeyId) > 0 {
		keyOutput, err := svc.GenerateDataKey(ctx, &awskms.GenerateDataKeyInput{
			KeyId:   &masterKeyId,
			KeySpec: awstypes.DataKeySpecAes256,
		})

		if err != nil {
			log.Error("data key generation error: " + err.Error())

			return nil, err
		}

		return keyOutput.CiphertextBlob, nil
	} else {
		return nil, nil
	}
}

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("aws configuration error: " + err.Error())
	}

	svc = awskms.NewFromConfig(cfg)
}

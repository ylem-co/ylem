package kms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"ylem_integrations/config"

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

func Encrypt(secret []byte, key []byte) ([]byte, error) {
	if len(secret) == 0 {
		return secret, nil
	}

	if len(key) == 0 {
		return secret, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(secret), nil)

	return ciphertext, nil
}

func Decrypt(secret []byte, key []byte) ([]byte, error) {
	if len(secret) == 0 {
		return secret, nil
	}

	if len(key) == 0 {
		return secret, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := secret[:nonceSize], secret[nonceSize:]

	plainsSecret, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainsSecret, nil
}

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("aws configuration error: " + err.Error())
	}

	svc = awskms.NewFromConfig(cfg)
}

func DecryptDataKey(ctx context.Context, masterKeyId string, dataKey []byte) ([]byte, error) {
	if len(masterKeyId) == 0 && len(dataKey) == 0 {
		return dataKey, nil
	}

	b, err := svc.Decrypt(ctx, &awskms.DecryptInput{
		CiphertextBlob:      dataKey,
		EncryptionAlgorithm: awstypes.EncryptionAlgorithmSpecSymmetricDefault,
		KeyId:               &masterKeyId,
	})

	if err != nil {
		return nil, err
	}

	return b.Plaintext, nil
}

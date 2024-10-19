package kms

import (
	"context"
	"ylem_taskrunner/config"

	messaging "github.com/ylem-co/shared-messaging"
)

func DecryptSource(s *messaging.SQLIntegration, ctx context.Context) error {
	var (
		value []byte
		err   error
	)

	keyId := config.Cfg().Aws.KmsKeyId
	decryptedDataKey, err := DecryptDataKey(ctx, keyId, s.DataKey)
	if err != nil {
		return err
	}

	if len(s.Host) > 0 {
		value, err = Decrypt(s.Host, decryptedDataKey)
		if err != nil {
			return err
		}

		s.Host = value
	}

	// Password
	if len(s.Password) > 0 {
		value, err = Decrypt(s.Password, decryptedDataKey)
		if err != nil {
			return err
		}

		s.Password = value
	}

	// SSH Host
	if len(s.SshHost) > 0 {
		value, err = Decrypt(s.SshHost, decryptedDataKey)
		if err != nil {
			return err
		}

		s.SshHost = value
	}

	// Credentials
	if len(s.Credentials) > 0 {
		value, err = Decrypt(s.Credentials, decryptedDataKey)
		if err != nil {
			return err
		}

		s.Credentials = value
	}

	return nil
}

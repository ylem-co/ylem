package entities

import (
	"context"
	"regexp"
	"ylem_integrations/config"
	"ylem_integrations/services/aws/kms"

	messaging "github.com/ylem-co/shared-messaging"
)

const IntegrationTypeSQL = "sql"

type SQLIntegration struct {
	Id               int64         `json:"-"`
	Integration      Integration   `json:"integration"`
	Type             string        `json:"type"`
	DataKey          kms.SecretBox `json:"-"`
	Host             kms.SecretBox `json:"host"`
	Port             int           `json:"port,omitempty"`
	User             string        `json:"user,omitempty"`
	Password         kms.SecretBox `json:"-"`
	Database         string        `json:"database,omitempty"`
	ConnectionType   string        `json:"connection_type"`
	SslEnabled       bool          `json:"ssl_enabled"`
	SshHost          kms.SecretBox `json:"ssh_host,omitempty"`
	SshPort          int           `json:"ssh_port,omitempty"`
	SshUser          string        `json:"ssh_user,omitempty"`
	ProjectId        *string       `json:"project_id,omitempty"`
	Credentials      kms.SecretBox `json:"credentials,omitempty"`
	EsVersion        *uint8        `json:"es_version,omitempty"`
	IsTrial          int           `json:"is_trial"`
}

const SQLIntegrationTypeMySQL = messaging.SQLIntegrationTypeMySQL
const SQLIntegrationTypeSnowflake = messaging.SQLIntegrationTypeSnowflake
const SQLIntegrationTypePostgresql = messaging.SQLIntegrationTypePostgresql
const SQLIntegrationTypeAWSRDS = messaging.SQLIntegrationTypeAWSRDS
const SQLIntegrationTypeGoogleCloudSQL = messaging.SQLIntegrationTypeGoogleCloudSQL
const SQLIntegrationTypeGoogleBigQuery = messaging.SQLIntegrationTypeGoogleBigQuery
const SQLIntegrationTypeElasticSearch = messaging.SQLIntegrationTypeElasticsearch
const SQLIntegrationTypePlanetScale = messaging.SQLIntegrationTypePlanetScale
const SQLIntegrationTypeImmuta = messaging.SQLIntegrationTypeImmuta
const SQLIntegrationTypeMicrosoftAzureSQL = messaging.SQLIntegrationTypeMicrosoftAzureSQL
const SQLIntegrationTypeRedshift = messaging.SQLIntegrationTypeRedshift
const SQLIntegrationTypeClickhouse = messaging.SQLIntegrationTypeClickhouse

const SQLIntegrationConnectionTypeDirect = messaging.SQLIntegrationConnectionTypeDirect
const SQLIntegrationConnectionTypeSsh = messaging.SQLIntegrationConnectionTypeSsh

func (s *SQLIntegration) IsDirectConnection() bool {
	return s.ConnectionType == SQLIntegrationConnectionTypeDirect
}

func (s *SQLIntegration) IsSshConnection() bool {
	return s.ConnectionType == SQLIntegrationConnectionTypeSsh
}

func (s *SQLIntegration) MaskHosts() {
	const trialHostConst = "Demo Host"

	if IsTrialHost(s.Host.PlainValue) {
		s.Host.SetPlainValue([]byte(trialHostConst))
	}

	if IsTrialHost(s.SshHost.PlainValue) {
		s.SshHost.SetPlainValue([]byte(trialHostConst))
	}
}

func (s *SQLIntegration) Decrypt(ctx context.Context) error {
	var (
		value []byte
		err   error
	)

	if s.DataKey.EncryptedValue == nil {
		s.DataKey.Open(s.DataKey.EncryptedValue)
		s.Host.Open(s.Host.EncryptedValue)
		s.Password.Open(s.Password.EncryptedValue)
		s.SshHost.Open(s.SshHost.EncryptedValue)
		s.Credentials.Open(s.Credentials.EncryptedValue)

		return nil
	}

	err = decryptDataKey(ctx, s)
	if err != nil {
		return err
	}

	// Host
	if s.Host.Sealed && s.Host.EncryptedValue != nil {
		value, err = kms.Decrypt(s.Host.EncryptedValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Host.Open(value)
	}

	// Password
	if s.Password.Sealed && s.Password.EncryptedValue != nil {
		value, err = kms.Decrypt(s.Password.EncryptedValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Password.Open(value)
	}

	// SSH Host
	if s.SshHost.Sealed && s.SshHost.EncryptedValue != nil {
		value, err = kms.Decrypt(s.SshHost.EncryptedValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.SshHost.Open(value)
	}

	// Credentials
	if s.Credentials.Sealed && s.Credentials.EncryptedValue != nil {
		value, err = kms.Decrypt(s.Credentials.EncryptedValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Credentials.Open(value)
	}

	return nil
}

func (s *SQLIntegration) Encrypt(ctx context.Context) error {
	var (
		value []byte
		err   error
	)

	err = decryptDataKey(ctx, s)
	if err != nil {
		return err
	}

	// Host
	if !s.Host.Sealed && s.Host.PlainValue != nil {
		value, err = kms.Encrypt(s.Host.PlainValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Host.SetEncryptedValue(value).Seal()
	}

	// Password
	if !s.Password.Sealed && s.Password.PlainValue != nil {
		value, err = kms.Encrypt(s.Password.PlainValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Password.SetEncryptedValue(value).Seal()
	}

	// SSH Host
	if !s.SshHost.Sealed && s.SshHost.PlainValue != nil {
		value, err = kms.Encrypt(s.SshHost.PlainValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.SshHost.SetEncryptedValue(value).Seal()
	}

	// Credentials
	if !s.Credentials.Sealed && s.Credentials.PlainValue != nil {
		value, err = kms.Encrypt(s.Credentials.PlainValue, s.DataKey.PlainValue)
		if err != nil {
			return err
		}
		s.Credentials.SetEncryptedValue(value).Seal()
	}

	return nil
}

type SQLIntegrationCollection struct {
	Items []SQLIntegration
}

func (s *SQLIntegrationCollection) MaskHosts() {
	const trialHostConst = "Demo Host"

	for k := range s.Items {
		if IsTrialHost(s.Items[k].Host.PlainValue) {
			s.Items[k].Host.SetPlainValue([]byte(trialHostConst))
		}

		if IsTrialHost(s.Items[k].SshHost.PlainValue) {
			s.Items[k].SshHost.SetPlainValue([]byte(trialHostConst))
		}
	}
}

func (s *SQLIntegrationCollection) Decrypt(ctx context.Context) error {
	for k := range s.Items {
		err := s.Items[k].Decrypt(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLIntegrationCollection) Encrypt(ctx context.Context) error {
	for k := range s.Items {
		err := s.Items[k].Encrypt(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func IsSQLIntegrationTypeSupported(Type string) bool {
	return map[string]bool{
		SQLIntegrationTypeMySQL:             true,
		SQLIntegrationTypeSnowflake:         true,
		SQLIntegrationTypePostgresql:        true,
		SQLIntegrationTypeAWSRDS:            true,
		SQLIntegrationTypeGoogleCloudSQL:    true,
		SQLIntegrationTypeGoogleBigQuery:    true,
		SQLIntegrationTypePlanetScale:       true,
		SQLIntegrationTypeImmuta:            true,
		SQLIntegrationTypeMicrosoftAzureSQL: true,
		SQLIntegrationTypeElasticSearch:     true,
		SQLIntegrationTypeRedshift:          true,
		SQLIntegrationTypeClickhouse:        true,
	}[Type]
}

func IsSQLIntegrationConnectionTypeSupported(Type string) bool {
	return Type == SQLIntegrationConnectionTypeDirect || Type == SQLIntegrationConnectionTypeSsh
}

func IsTrialHost(host []byte) bool {
	matched, err := regexp.MatchString(`crnutbhzybdn|ec2-35-158-203-179`, string(host))

	if err != nil {
		return false
	}

	return matched
}

func decryptDataKey(ctx context.Context, s *SQLIntegration) error {
	if s.DataKey.Sealed {
		decrypted, err := kms.DecryptDataKey(ctx, config.Cfg().Aws.KmsKeyId, s.DataKey.EncryptedValue)
		if err != nil {
			return err
		}

		s.DataKey.Open(decrypted)
	}

	return nil
}

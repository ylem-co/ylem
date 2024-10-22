package repositories

import (
	"database/sql"
	"ylem_integrations/entities"
	"ylem_integrations/services/aws/kms"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateSQLIntegration(db *sql.DB, entity *entities.SQLIntegration) (*entities.SQLIntegration, error) {
	log.Tracef("Creating SQL Integration")

	tx, err := db.Begin()
	defer tx.Rollback() //nolint:all
	if err != nil {
		log.Error(err)

		return nil, err
	}

	entity.Integration.Uuid = uuid.NewString()
	entity.Integration.Status = entities.IntegrationStatusNew
	entity.Integration.Type = entities.IntegrationTypeSQL
	entity.Integration.IoType = entities.IntegrationIoTypeReadWrite

	if entity.Type == entities.SQLIntegrationTypeElasticSearch {
		entity.Integration.IoType = entities.IntegrationIoTypeRead
	}

	{
		Query := `INSERT INTO integrations 
        (uuid, creator_uuid, organization_uuid, name, status, value, type, io_type, user_updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return nil, err
		}
		defer stmt.Close()

		result, err := stmt.Exec(
			entity.Integration.Uuid,
			entity.Integration.CreatorUuid,
			entity.Integration.OrganizationUuid,
			entity.Integration.Name,
			entity.Integration.Status,
			entity.Integration.Value,
			entity.Integration.Type,
			entity.Integration.IoType,
			entity.Integration.UserUpdatedAt,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return nil, err
		}

		entity.Integration.Id, _ = result.LastInsertId()
	}

	{
		Query := "INSERT INTO sqls (`integration_id`, `data_key`, `host`, `port`, `user`, `password`, `database`, `connection_type`, `ssl_enabled`, `ssh_host`, `ssh_port`, `ssh_user`, `project_id`, `credentials`, `es_version`, `is_trial`, `type`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		stmt, err := tx.Prepare(Query)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return nil, err
		}
		defer stmt.Close()

		sSslEnabled := 0
		if entity.SslEnabled {
			sSslEnabled = 1
		}

		_, err = stmt.Exec(
			entity.Integration.Id,
			entity.DataKey.EncryptedValue,
			entity.Host.EncryptedValue,
			entity.Port,
			entity.User,
			entity.Password.EncryptedValue,
			entity.Database,
			entity.ConnectionType,
			sSslEnabled,
			entity.SshHost.EncryptedValue,
			entity.SshPort,
			entity.SshUser,
			entity.ProjectId,
			entity.Credentials.EncryptedValue,
			entity.EsVersion,
			entity.IsTrial,
			entity.Type,
		)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err)
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return entity, nil
}

func FindSQLIntegration(db *sql.DB, Uuid string) (*entities.SQLIntegration, error) {
	log.Trace("Finding SQL Integration")

	Query := `SELECT
		d.id as integration_id,
		d.uuid,
		d.creator_uuid,
		d.organization_uuid,
		d.status,
		d.type,
		d.io_type,
		d.name,
		d.value,
		d.user_updated_at,
		s.id,
		s.type,
		s.data_key,
		s.host,
		s.port,
		s.user,
		s.password,
		s.database,
		s.connection_type,
		s.ssl_enabled,
		s.ssh_host,
		s.ssh_port,
		s.ssh_user,
		s.project_id,
		IFNULL(s.credentials, ""),
		s.es_version,
		s.is_trial
	FROM
		sqls s
		INNER JOIN integrations d ON d.id = s.integration_id
	WHERE 
		d.uuid = ?
		AND d.deleted_at IS NULL`

	var entity entities.SQLIntegration

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)

		return nil, err
	}
	defer stmt.Close()

	var (
		dataKey     []byte
		host        []byte
		password    []byte
		sshHost     []byte
		credentials []byte
	)

	err = stmt.QueryRow(Uuid).Scan(
		&entity.Integration.Id,
		&entity.Integration.Uuid,
		&entity.Integration.CreatorUuid,
		&entity.Integration.OrganizationUuid,
		&entity.Integration.Status,
		&entity.Integration.Type,
		&entity.Integration.IoType,
		&entity.Integration.Name,
		&entity.Integration.Value,
		&entity.Integration.UserUpdatedAt,
		&entity.Id,
		&entity.Type,
		&dataKey,
		&host,
		&entity.Port,
		&entity.User,
		&password,
		&entity.Database,
		&entity.ConnectionType,
		&entity.SslEnabled,
		&sshHost,
		&entity.SshPort,
		&entity.SshUser,
		&entity.ProjectId,
		&credentials,
		&entity.EsVersion,
		&entity.IsTrial,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	entity.DataKey = kms.NewSealedSecretBox(dataKey)
	entity.Host = kms.NewSealedSecretBox(host)
	entity.Password = kms.NewSealedSecretBox(password)
	entity.SshHost = kms.NewSealedSecretBox(sshHost)
	entity.Credentials = kms.NewSealedSecretBox(credentials)

	return &entity, nil
}

func UpdateSQLIntegration(db *sql.DB, entity *entities.SQLIntegration, rewritePassword bool) error {
	Query := `UPDATE 
		sqls s
		INNER JOIN integrations d ON d.id = s.integration_id
	       SET d.name = ?,
	       s.type = ?,
	       s.ssl_enabled = ?,
	       d.status = ?,
	       s.host = ?,
	       s.port = ?,
	       s.user = ?,
	       s.database = ?,
	       s.connection_type = ?,
	       s.ssh_host = ?,
	       s.ssh_port = ?,
	       s.ssh_user = ?,
	       s.project_id = ?,
	       s.credentials = ?,
	       s.es_version = ?,
		   d.user_updated_at = ?
	       WHERE d.id = ?
	       `

	if rewritePassword {
		Query = `UPDATE 
		sqls s
		INNER JOIN integrations d ON d.id = s.integration_id
	       SET d.name = ?,
	       s.type = ?,
	       s.ssl_enabled = ?,
	       d.status = ?,
	       s.host = ?,
	       s.port = ?,
	       s.user = ?,
	       s.database = ?,
	       s.connection_type = ?,
	       s.ssh_host = ?,
	       s.ssh_port = ?,
	       s.ssh_user = ?,
	       s.project_id = ?,
	       s.credentials = ?,
	       s.es_version = ?,
	       s.password = ?,
		   d.user_updated_at = ?
	       WHERE d.id = ?
	       `
	}

	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())

		return err
	}

	sSslEnabled := 0
	if entity.SslEnabled {
		sSslEnabled = 1
	}

	if rewritePassword {
		_, err = stmt.Exec(
			entity.Integration.Name,
			entity.Type,
			sSslEnabled,
			entity.Integration.Status,
			entity.Host.EncryptedValue,
			entity.Port,
			entity.User,
			entity.Database,
			entity.ConnectionType,
			entity.SshHost.EncryptedValue,
			entity.SshPort,
			entity.SshUser,
			entity.ProjectId,
			entity.Credentials.EncryptedValue,
			entity.EsVersion,
			entity.Integration.UserUpdatedAt,
			entity.Password.EncryptedValue,
			entity.Integration.Id,
		)
	} else {
		_, err = stmt.Exec(
			entity.Integration.Name,
			entity.Type,
			sSslEnabled,
			entity.Integration.Status,
			entity.Host.EncryptedValue,
			entity.Port,
			entity.User,
			entity.Database,
			entity.ConnectionType,
			entity.SshHost.EncryptedValue,
			entity.SshPort,
			entity.SshUser,
			entity.ProjectId,
			entity.Credentials.EncryptedValue,
			entity.EsVersion,
			entity.Integration.UserUpdatedAt,
			entity.Integration.Id,
		)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

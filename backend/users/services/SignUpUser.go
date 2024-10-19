package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"ylem_users/entities"
	"ylem_users/repositories"
	"ylem_users/services/kms"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	log "github.com/sirupsen/logrus"
)

const (
	PublicEmailGmail   = "gmail.com"
	PublicEmailYahoo   = "yahoo.com"
	PublicEmailHotmail = "hotmail.com"
)

type SignUpUser struct {
	db  *sql.DB
	ctx context.Context
}

func (s *SignUpUser) SignUpExternalUser(user goth.User, organisationName string, source string, externalId string, orgId int64) (*entities.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Error(err)

		return nil, err
	}
	defer tx.Rollback() //nolint:all

	userId, err := s.createUser(tx, user, source, externalId)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return nil, err
	}

	if orgId == 0 {
		orgId, err = s.createOrganisation(tx, organisationName, userId)
		if err != nil {
			_ = tx.Rollback()
			log.Error(err)

			return nil, err
		}
	}

	err = s.assignOrganizationToUser(tx, userId, orgId)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return nil, err
	}

	return repositories.GetUserByExternalSystemId(s.db, source, externalId)
}

func (s *SignUpUser) createUser(tx *sql.Tx, user goth.User, source string, externalId string) (int64, error) {
	query := `INSERT INTO users 
        (first_name, last_name, email, password, phone, uuid, source, external_system_id, roles, email_confirmation_token, is_email_confirmed) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return 0, err
	}
	defer stmt.Close()

	roles, _ := json.Marshal([]string{entities.ROLE_ORGANIZATION_ADMIN})
	userUuid := uuid.NewString()
	result, err := stmt.Exec(user.FirstName, user.LastName, user.Email, "", "", userUuid, source, externalId, roles, "", 1)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return 0, err
	}

	return result.LastInsertId()
}

func (s *SignUpUser) createOrganisation(tx *sql.Tx, name string, ownerId int64) (int64, error) {
	query := `INSERT INTO organizations 
        (name, data_key, creator_id, uuid) 
        VALUES (?, ?, ?, ?)
        `

    dataKey, err := kms.IssueDataKeyWithContext(s.ctx)
    if err != nil {
		log.Error(err)

		return 0, nil
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return 0, nil
	}
	defer stmt.Close()

	orgUuid := uuid.NewString()
	res, err := stmt.Exec(name, dataKey, ownerId, orgUuid)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err.Error())

		return 0, nil
	}

	return res.LastInsertId()
}

func (s *SignUpUser) assignOrganizationToUser(tx *sql.Tx, userId int64, orgId int64) error {
	query := `UPDATE users 
        SET organization_id = ? 
        WHERE id = ?
        `

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(orgId, userId)
	if err != nil {
		_ = tx.Rollback()
		log.Error(err)

		return err
	}

	return nil
}

func CreateSignUpSrv(db *sql.DB, ctx context.Context) SignUpUser {
	return SignUpUser{
		db:  db,
		ctx: ctx,
	}
}

func IsGenericEmailProvider(host string) bool {
	switch host {
	case
		PublicEmailGmail,
		PublicEmailYahoo,
		PublicEmailHotmail:

		return true
	}

	return false
}

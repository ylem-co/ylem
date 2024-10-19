package api

import (
	"context"
	"strings"
	"database/sql"
	"encoding/json"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"
	"ylem_users/services/kms"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	log "github.com/sirupsen/logrus"
)

const emailConfirmationTokenLength = 24

func (auth *AuthMiddleware) PostUser(w http.ResponseWriter, r *http.Request) {
	var user services.HttpUser
	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &user)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateHttpUser(user, w)
	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	/*if repositories.DoesOrganizationExist(db, user.OrganizationName, "") {
		rp, _ := json.Marshal(map[string]string{"error": "Organization already exist", "fields": "organization_exists"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}*/

	if repositories.DoesUserExist(db, user.Email, "") {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid Email", "fields": "email"})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}
	var (
		ok         bool
	)
	if user.InvitationKey != "" {
		ok, _ = SaveInvitedUser(db, user)
	} else {
		ok, _ = SaveOrganizationAndUser(db, user, auth)
	}

	if ok {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}

func SaveInvitedUser(db *sql.DB, user services.HttpUser) (bool, *entities.User) {
	log.Tracef("save invite user")

	// Get a Tx for making transaction requests.
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	defer tx.Rollback() //nolint:all

	/* Getting organization ID by invitation key */

	var orgId int
	OrganizationQuery := `SELECT IF(organization_id IS NULL, 0, organization_id) as id
              FROM invitations
              WHERE invitation_code = ? AND accepter_id IS NULL
              `
	stmt, _ := tx.Prepare(OrganizationQuery)
	err = stmt.QueryRow(user.InvitationKey).Scan(&orgId)

	if err != nil || orgId == 0 {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	/* Creating user */

	insertUserQuery := `INSERT INTO users 
        (first_name, last_name, email, phone, uuid, password, roles, organization_id, email_confirmation_token, is_email_confirmed) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
        `

	insertStatement, err := tx.Prepare(insertUserQuery)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}
	defer insertStatement.Close()

	hash, _ := HashPassword(user.Password)
	roles, _ := json.Marshal([]string{entities.ROLE_TEAM_MEMBER})
	userUuid := uuid.NewString()
	emailConfirmationToken := helpers.CreateRandomNumericString(emailConfirmationTokenLength)
	_, err = insertStatement.Exec(user.FirstName, user.LastName, user.Email, user.Phone, userUuid, hash, roles, orgId, emailConfirmationToken)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	/* Getting user ID */

	var usr entities.User
	UserQuery := `SELECT id, first_name, last_name, email, phone, uuid, roles, is_email_confirmed
              FROM users
              WHERE uuid = ?
              `
	stmt, _ = tx.Prepare(UserQuery)
	err = stmt.QueryRow(userUuid).Scan(&usr.Id, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Phone, &usr.Uuid, &usr.Roles, &usr.IsEmailConfirmed)

	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	/* Updating invitation key */

	updateOrgQuery := `UPDATE invitations 
        SET accepter_id = ? 
        WHERE invitation_code = ?
        `

	updateStatement, err := tx.Prepare(updateOrgQuery)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}
	defer updateStatement.Close()

	_, err = updateStatement.Exec(usr.Id, user.InvitationKey)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	return true, &usr
}

func SaveOrganizationAndUser(db *sql.DB, user services.HttpUser, auth *AuthMiddleware) (bool, *entities.User) {
	log.Tracef("save organization and user")

	return SaveOrganizationAndUserWithContext(context.Background(), db, user, auth)
}

func SaveOrganizationAndUserWithContext(ctx context.Context, db *sql.DB, user services.HttpUser, auth *AuthMiddleware) (bool, *entities.User) {
	log.Tracef("save organization and user with context")

	// Get a Tx for making transaction requests.
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return false, nil
	}

	defer tx.Rollback() //nolint:all

	org, orgEx := repositories.GetOrganization(db)
	roles, _ := json.Marshal([]string{entities.ROLE_ORGANIZATION_ADMIN})
	orgUuid := uuid.NewString()
	if !orgEx {

		/* Creating organization */
		insertOrgQuery := `INSERT INTO organizations 
	        (name, data_key, uuid) 
	        VALUES (?, ?, ?)
	        `

	    dataKey, err := kms.IssueDataKeyWithContext(ctx)

		if err != nil {
			_ = tx.Rollback()
			log.Error(err.Error())
			return false, nil
		}

		insertStatement, err := tx.Prepare(insertOrgQuery)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return false, nil
		}
		defer insertStatement.Close()

		_, err = insertStatement.Exec("My organization", dataKey, orgUuid)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return false, nil
		}
	} else {
		roles, _ = json.Marshal([]string{entities.ROLE_TEAM_MEMBER})
		orgUuid = org.Uuid
	}

	/* Creating user */

	insertUserQuery := `INSERT INTO users 
        (first_name, last_name, email, phone, uuid, password, roles, email_confirmation_token, is_email_confirmed) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `

	insertStatement, err := tx.Prepare(insertUserQuery)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}
	defer insertStatement.Close()

	hash, _ := HashPassword(user.Password)
	userUuid := uuid.NewString()
	emailConfirmationToken := helpers.CreateRandomNumericString(emailConfirmationTokenLength)
	_, err = insertStatement.Exec(user.FirstName, user.LastName, user.Email, user.Phone, userUuid, hash, roles, emailConfirmationToken, 1)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	var orgId int
	var usr entities.User
	UserQuery := `SELECT id, first_name, last_name, email, phone, uuid, roles
	              FROM users
	              WHERE uuid = ?
	              `
	stmt, _ := tx.Prepare(UserQuery)
	err = stmt.QueryRow(userUuid).Scan(&usr.Id, &usr.FirstName, &usr.LastName, &usr.Email, &usr.Phone, &usr.Uuid, &usr.Roles)

	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, nil
	}

	if !orgEx {
		var organization entities.Organization
		OrganizationQuery := `SELECT id, name
	              FROM organizations
	              WHERE uuid = ?
	              `
		stmt, _ = tx.Prepare(OrganizationQuery)
		err = stmt.QueryRow(orgUuid).Scan(&organization.Id, &organization.Name)

		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return false, nil
		}

		/* Updating organization */

		updateOrgQuery := `UPDATE organizations 
	        SET creator_id = ? 
	        WHERE uuid = ?
	        `

		updateStatement, err := tx.Prepare(updateOrgQuery)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return false, &usr
		}
		defer updateStatement.Close()

		_, err = updateStatement.Exec(usr.Id, orgUuid)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err.Error())
			return false, &usr
		}

		orgId = organization.Id
	} else {
		orgId = org.Id
	}

	/* Updating user */

	updateUserQuery := `UPDATE users 
        SET organization_id = ? 
        WHERE uuid = ?
        `

	updateStatement, err := tx.Prepare(updateUserQuery)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, &usr
	}
	defer updateStatement.Close()

	_, err = updateStatement.Exec(orgId, userUuid)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, &usr
	}

	tokenErr, tokenString := CreateJWTToken(usr, auth, nil)
	if tokenErr != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, &usr
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Println(err.Error())
		return false, &usr
	}

	if !orgEx {
		log.Tracef("create default email destination")
		dUuid, ok := services.CreateDefaultEmailDestination(orgUuid, user.Email, user.FirstName, tokenString)

		if ok {
			log.Tracef("create trial db data source")
			sUuid, ok := services.CreateTrialDBDataSource(orgUuid, tokenString)
			if ok {
				log.Tracef("test trial database connection")
				services.TestTrialDBDataSource(sUuid, tokenString)
				log.Tracef("create trial pipelines")
				_ = services.CreateTrialPipelines(orgUuid, dUuid, sUuid, tokenString)
			}
		}
	}

	return true, &usr
}

func ValidateHttpUser(user services.HttpUser, w http.ResponseWriter) []string {
	var errorFields []string

	if user.FirstName == "" {
		errorFields = append(errorFields, "first_name")
	}

	if user.LastName == "" {
		errorFields = append(errorFields, "last_name")
	}

	/*if user.OrganizationName == "" && user.InvitationKey == "" {
		errorFields = append(errorFields, "organization_name")
		errorFields = append(errorFields, "invitation_key")
	}*/

	if user.Email == "" || !govalidator.IsEmail(user.Email) {
		errorFields = append(errorFields, "email")
	}

	if user.Password == "" || !services.IsPasswordValid(user.Password) {
		errorFields = append(errorFields, "password")
	}

	if user.Password != user.ConirmPassword {
		errorFields = append(errorFields, "password")
		errorFields = append(errorFields, "confirm_password")
	}

	if user.Phone != "" {
		if !services.IsPhoneValid(user.Phone) {
			errorFields = append(errorFields, "phone")
		}
	}

	return errorFields
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

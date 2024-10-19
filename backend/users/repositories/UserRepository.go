package repositories

import (
	"database/sql"
	"ylem_users/entities"
	"ylem_users/helpers"

	log "github.com/sirupsen/logrus"
)

func DoesUserExist(db *sql.DB, email string, userUuid string) bool {
	Query := ""
	var rows *sql.Rows

	if userUuid == "" {
		Query = `SELECT COUNT(*)
              FROM users
              WHERE email = ?
              `
		rows, _ = db.Query(Query, email)
	} else {
		Query = `SELECT COUNT(*)
              FROM users
              WHERE email = ? AND uuid != ?
              `
		rows, _ = db.Query(Query, email, userUuid)
	}

	if helpers.NumRows(rows) > 0 {
		return true
	}

	return false
}

func GetUserByEmail(db *sql.DB, email string) (entities.User, bool) {
	Query := `SELECT id, first_name, last_name, uuid, email, password, phone, roles, is_email_confirmed
              FROM users
              WHERE email = ? AND is_active = 1
              `
	var usr entities.User
	stmt, err := db.Prepare(Query)
	if err != nil {
		log.Error(err)
		return usr, false
	}
	err = stmt.QueryRow(email).Scan(&usr.Id, &usr.FirstName, &usr.LastName, &usr.Uuid, &usr.Email, &usr.HashedPassword, &usr.Phone, &usr.Roles, &usr.IsEmailConfirmed)

	if err != nil {
		log.Error(err)
		return usr, false
	} else {
		return usr, true
	}
}

func GetUserByExternalSystemId(db *sql.DB, source string, externalSystemId string) (*entities.User, error) {
	query := `SELECT id, first_name, last_name, uuid, email, password, phone, roles, is_email_confirmed
              FROM users
              WHERE source = ? 
                AND external_system_id = ?
                AND is_active = 1
                AND is_email_confirmed = 1
              `

	var usr entities.User
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRow(source, externalSystemId).
		Scan(
			&usr.Id,
			&usr.FirstName,
			&usr.LastName,
			&usr.Uuid,
			&usr.Email,
			&usr.HashedPassword,
			&usr.Phone,
			&usr.Roles,
			&usr.IsEmailConfirmed,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func GetUserByUuid(db *sql.DB, uuid string) (entities.User, bool) {
		Query := `SELECT id, first_name, last_name, uuid, email, password, roles, organization_id, is_email_confirmed
	              FROM users
	              WHERE uuid = ? AND is_active = 1
	              `
		var usr entities.User
		stmt, err := db.Prepare(Query)
		if err != nil {
			log.Error(err)
			return usr, false
		}
		err = stmt.QueryRow(uuid).Scan(&usr.Id, &usr.FirstName, &usr.LastName, &usr.Uuid, &usr.Email, &usr.HashedPassword, &usr.Roles, &usr.OrganizationId, &usr.IsEmailConfirmed)

		if err != nil {
			log.Error(err)
			return usr, false
		} else {
			return usr, true
		}
}

func GetUsersByOrganizationUuid(db *sql.DB, uuid string) ([]entities.UserToExpose, bool) {
		Query := `SELECT u.first_name, u.last_name, u.email, u.uuid, u.roles, u.is_active, u.is_email_confirmed
	              FROM users u
	              JOIN organizations o ON o.id = u.organization_id
	              WHERE o.uuid = ?
	              ORDER BY is_active DESC, u.created_at DESC
	              `
		stmt, err := db.Prepare(Query)
		if err != nil {
			log.Error(err)
		}
		rows, err := stmt.Query(uuid)
		if err != nil {
			log.Error(err)
			return nil, false
		}
		defer rows.Close()

		var users []entities.UserToExpose
		var firstName string
		var lastName string
		var email string
		var userUuid string
		var roles string
		var isActive int
		var isEmailConfirmed int

		for rows.Next() {
			err := rows.Scan(&firstName, &lastName, &email, &userUuid, &roles, &isActive, &isEmailConfirmed)
			if err != nil {
				return nil, false
			}

			users = append(users, entities.UserToExpose{FirstName: firstName, LastName: lastName, Email: email, Uuid: userUuid, Roles: roles, IsActive: isActive, IsEmailConfirmed: isEmailConfirmed})
		}
		if err := rows.Err(); err != nil {
			return users, false
		}

		return users, true
}

func DeleteUser(db *sql.DB, uuid string) bool {
		updateQuery := `UPDATE users 
	        SET is_active = 0
	        WHERE uuid = ?
	        `

		updateStatement, err := db.Prepare(updateQuery)
		if err != nil {
			return false
		}
		defer updateStatement.Close()

		_, err = updateStatement.Exec(uuid)
		if err != nil {
			log.Error(err)
			return false
		}

		return true
}

func ActivateUser(db *sql.DB, uuid string) bool {
		updateQuery := `UPDATE users 
	        SET is_active = 1
	        WHERE uuid = ?
	        `

		updateStatement, err := db.Prepare(updateQuery)
		if err != nil {
			return false
		}
		defer updateStatement.Close()

		_, err = updateStatement.Exec(uuid)
		if err != nil {
			log.Error(err)
			return false
		}

		return true
}

func AssignRole(db *sql.DB, uuid string, role []byte) bool {
		updateQuery := `UPDATE users 
	        SET roles = ?
	        WHERE uuid = ?
	        `

		updateStatement, err := db.Prepare(updateQuery)
		if err != nil {
			return false
		}
		defer updateStatement.Close()

		_, err = updateStatement.Exec(role, uuid)
		if err != nil {
			log.Error(err)
			return false
		}

		return true
}

func GetUserByEmailToken(db *sql.DB, emailToken string) (*entities.User, error) {
		Query := `SELECT uuid, is_email_confirmed
	              FROM users
	              WHERE email_confirmation_token = ? AND is_active = 1
	              `
		var usr entities.User
		stmt, err := db.Prepare(Query)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = stmt.QueryRow(emailToken).Scan(&usr.Uuid, &usr.IsEmailConfirmed)

		if err == sql.ErrNoRows {
			return nil, nil
		}

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		return &usr, nil
}

func UpdateUserIsEmailConfirmed(db *sql.DB, usr *entities.User, isEmailConfirmed int) error {
		updateOrgQuery := `UPDATE users 
	        SET is_email_confirmed = ? 
	        WHERE uuid = ?
	        `

		stmt, err := db.Prepare(updateOrgQuery)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(isEmailConfirmed, usr.Uuid)

		if err != nil && err != sql.ErrNoRows {
			log.Println(err.Error())
			return err
		}
		return nil
}

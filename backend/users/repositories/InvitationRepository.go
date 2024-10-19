package repositories

import (
	"database/sql"
	"ylem_users/entities"
)

func GetPendingInvitationsByOrganizationUuid(db *sql.DB, uuid string) ([]entities.InvitationToExpose, bool) {
	Query := `SELECT i.uuid, i.email, i.created_at, i.invitation_code
              FROM invitations i
              JOIN organizations o ON o.id = i.organization_id
              WHERE o.uuid = ? AND i.accepter_id IS NULL
              ORDER BY i.created_at DESC
              `
	stmt, err := db.Prepare(Query)
	if err != nil {
		return nil, false
	}

	rows, err := stmt.Query(uuid)
	if err != nil {
		return nil, false
	}
	defer rows.Close()

	var invitations []entities.InvitationToExpose
	var email string
	var invitationUuid string
	var createdAt string
	var invitationCode string

	for rows.Next() {
		err := rows.Scan(&invitationUuid, &email, &createdAt, &invitationCode)
		if err != nil {
			return nil, false
		}

		invitations = append(invitations, entities.InvitationToExpose{Uuid: invitationUuid, Email: email, CreatedAt: createdAt, InvitationCode: invitationCode})
	}
	if err := rows.Err(); err != nil {
		return invitations, false
	}

	return invitations, true
}

func SaveInvitation(db *sql.DB, uuid string, email string, invitationCode string, orgId int, userId int) bool {
	insertQuery := `INSERT INTO invitations 
        (uuid, email, invitation_code, organization_id, sender_id) 
        VALUES (?, ?, ?, ?, ?)
        `

	insertStatement, err := db.Prepare(insertQuery)
	if err != nil {
		return false
	}
	defer insertStatement.Close()

	_, err = insertStatement.Exec(uuid, email, invitationCode, orgId, userId)
	return err == nil
}

func ValidateInvitationByKey(db *sql.DB, key string) bool {
	Query := `SELECT IF(id IS NULL, 0, 1) as status
              FROM invitations
              WHERE invitation_code = ? AND accepter_id IS NULL
              `
	var id int
	stmt, err := db.Prepare(Query)
	if err != nil {
		return false
	}

	err = stmt.QueryRow(key).Scan(&id)

	if err != nil || id == 0 {
		return false
	} else {
		return true
	}
}

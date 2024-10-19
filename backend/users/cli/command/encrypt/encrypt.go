package encrypt

import (
	"time"
	"ylem_users/helpers"
	"ylem_users/services/kms"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var issueDataKey cli.ActionFunc = func(c *cli.Context) error {
	db := helpers.DbConn()

	query := "SELECT `uuid` FROM organizations WHERE data_key IS NULL"

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Error(err.Error())

		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Error(err)

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			uuid                 string
			encryptedDataKey     []byte
		)

		err := rows.Scan(
			&uuid,
		)

		if err != nil {
			log.Error(err.Error())

			continue
		}

		log.Debugf("Issue data key for the organization uuid %s", uuid)
		encryptedDataKey, err = kms.IssueDataKeyWithContext(c.Context)
		if err != nil {
			log.Error(err.Error())

			return err
		}

		// update

		updateQuery := "UPDATE organizations SET `data_key` = ? WHERE uuid = ?"
		stmt, err = db.Prepare(updateQuery)
		if err != nil {
			log.Error(err.Error())

			return err
		}

		defer stmt.Close()
		_, err = stmt.Exec(
			encryptedDataKey,
			uuid,
		)

		if err != nil {
			log.Error(err.Error())

			return err
		}

		log.Infof("Successfully issued data key for  the organization %s", uuid)

		time.Sleep(1 * time.Second)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}

var issueDataKeys = &cli.Command{
	Name:   "issue",
	Usage:  "Issue data keys for organizations",
	Action: issueDataKey,
}

var Command = &cli.Command{
	Name:  "encrypt",
	Usage: "Encrypt database values",
	Subcommands: []*cli.Command{
		issueDataKeys,
	},
}

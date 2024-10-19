package command

import (
	"ylem_users/helpers"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var migrateHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Applying migrations to database")

	driver, err := mysql.WithInstance(helpers.DbConn(), &mysql.Config{})
	if err != nil {
		log.Info("migrate command failed: " + err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"mysql",
		driver,
	)
	if err != nil {
		log.Info("migrate command failed: " + err.Error())
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Info("Nothing to migrate")
		return nil
	}

	return err
}

var MigrateCommand = &cli.Command{
	Name:   "migrate",
	Usage:  "Migrate the DB to the latest version",
	Action: migrateHandler,
}

var MigrationsCommand = &cli.Command{
	Name:  "migrations",
	Usage: "Migration-related commands",
	Subcommands: []*cli.Command{
		MigrateCommand,
	},
}

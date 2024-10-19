package command

import (
	"ylem_users/helpers"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var fixtureLoadHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Info("Loading fixtures...")
	db := helpers.DbConn()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Info("Done.")

	return nil
}

var FixtureLoadCommand = &cli.Command{
	Name:   "load",
	Usage:  "Load fixtures into database",
	Action: fixtureLoadHandler,
}

var FixturesCommand = &cli.Command{
	Name:  "fixtures",
	Usage: "Fixtures",
	Subcommands: []*cli.Command{
		FixtureLoadCommand,
	},
}

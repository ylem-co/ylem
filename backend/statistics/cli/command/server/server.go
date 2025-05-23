package server

import (
	"ylem_statistics/config"
	"ylem_statistics/services/server"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var serveHandler cli.ActionFunc = func(c *cli.Context) error {
	log.Debug("serve command called")
	return server.NewServer(c.Context, config.Cfg().Listen).Run()
}

var ServeCommand = &cli.Command{
	Name:        "serve",
	Description: "Start a HTTP(s) server",
	Usage:       "Start a HTTP(s) server",
	Action:      serveHandler,
}

var Command = &cli.Command{
	Name:        "server",
	Description: "HTTP(s) server commands",
	Usage:       "HTTP(s) server commands",
	Subcommands: []*cli.Command{
		ServeCommand,
	},
}

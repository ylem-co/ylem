package main

import (
	"context"
	"os"
	"time"
	"ylem_taskrunner/cli"

	log "github.com/sirupsen/logrus"
)

func main() {
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := cli.NewApplication().RunContext(ctx, os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"os"
	"os/signal"
	"ylem_api/cli"
	"syscall"
	"time"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func main() {
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	govalidator.SetFieldsRequiredByDefault(true)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool)
	go func() {
		defer close(done)
		err := cli.NewApplication().RunContext(ctx, os.Args)
		if err != nil {
			log.Fatalf("error running application: %v", err)
		} else {
			log.Info("Graceful shutdown")
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-wait: // wait for SIGINT/SIGTERM
		signal.Reset(syscall.SIGINT, syscall.SIGTERM) // resetting signal listener, so that repeated Ctrl+C will exit immediately
		cancel()                                      // graceful stop
		<-done

	case <-done:
		cancel() // graceful stop
	}
}

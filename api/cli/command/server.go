package command

import (
	"context"
	"net"
	"net/http"
	"ylem_api/api"
	"ylem_api/api/stats"
	"ylem_api/api/pipeline"
	"ylem_api/config"
	"ylem_api/service/oauth"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var serveHandler cli.ActionFunc = func(c *cli.Context) error {
	rtr := mux.NewRouter()

	internal := rtr.PathPrefix("/oauth-api").Subrouter()
	internal.HandleFunc("/clients/", api.ListClientsByUser).Methods(http.MethodGet)
	internal.HandleFunc("/clients/", api.CreateOauthClient).Methods(http.MethodPost)
	internal.HandleFunc("/client/{uuid}/delete/", api.DeleteOauthClient).Methods(http.MethodPost)

	externalV1 := rtr.PathPrefix("/v1").Subrouter()
	externalV1.HandleFunc("/oauth/token", api.GenerateToken).Methods(http.MethodGet, http.MethodPost)
	externalV1.HandleFunc("/pipelines/{pipelineUuid}/runs/", api.AuthenticateScoped(pipeline.RunPipeline, oauth.ScopePipelinesRun)).Methods(http.MethodPost)

	externalV1.HandleFunc("/stats/tasks/{uuid}/stats/{dateFrom}/{dateTo}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/tasks/{uuid}/aggregated-stats/{dateFrom}/{period}/{periodCount}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/tasks/{uuid}/last-run/stats", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)

	externalV1.HandleFunc("/stats/pipelines/{uuid}/stats/{dateFrom}/{dateTo}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/pipelines/{uuid}/aggregated-stats/{dateFrom}/{period}/{periodCount}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/pipelines/{uuid}/last-run/stats", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/pipelines/{uuid}/last-runs-log/{dateFrom}/{dateTo}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)
	externalV1.HandleFunc("/stats/pipelines/{uuid}/logs/{dateFrom}/{dateTo}", api.AuthenticateScoped(stats.ProxyRequest, oauth.ScopeStatsRead)).Methods(http.MethodGet)

	http.Handle("/", rtr)

	log.Infof("Listening on %s", config.Cfg().Listen)

	server := &http.Server{
		Addr:    config.Cfg().Listen,
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return c.Context
		},
	}

	go func() {
		<-c.Done()
		_ = server.Shutdown(c.Context)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

var ServeCommand = &cli.Command{
	Name:        "serve",
	Description: "Start a HTTP(s) server",
	Usage:       "Start a HTTP(s) server",
	Action:      serveHandler,
}

var ServerCommands = &cli.Command{
	Name:        "server",
	Description: "HTTP(s) server commands",
	Usage:       "HTTP(s) server commands",
	Subcommands: []*cli.Command{
		ServeCommand,
	},
}

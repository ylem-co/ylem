package server

import (
	"context"
	"net"
	"net/http"
	"ylem_statistics/api"
	"ylem_statistics/config"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Listen string
	ctx    context.Context
}

func (s *Server) Run() error {
	log.Info("Starting server listening on " + s.Listen)

	rtr := mux.NewRouter()

	rtr.HandleFunc("/organization/{uuid}/tasks/slow/{dateFrom}/{dateTo}/{threshold}/{type}", api.GetSlowTaskRuns).Methods(http.MethodGet)
	rtr.HandleFunc("/tasks/{uuid}/stats/{dateFrom}/{dateTo}", api.NewGetStatsFunc(api.TaskStatsLoadStatsFunc)).Methods(http.MethodGet)
	rtr.HandleFunc("/tasks/{uuid}/aggregated-stats/{dateFrom}/{period}/{periodCount}", api.NewGetAggregatedStats(api.TaskLoadAggregatedStatsFunc)).Methods(http.MethodGet)
	rtr.HandleFunc("/tasks/{uuid}/last-run/stats", api.NewGetLastRunStats(api.TaskLastRunStatsLoadStateFunc)).Methods(http.MethodGet)
	
	rtr.HandleFunc("/pipelines/{uuid}/stats/{dateFrom}/{dateTo}", api.NewGetStatsFunc(api.PipelineStatsLoadStatsFunc)).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/aggregated-stats/{dateFrom}/{period}/{periodCount}", api.NewGetAggregatedStats(api.PipelineLoadAggregatedStatsFunc)).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/last-run/stats", api.NewGetLastRunStats(api.PipelineLastRunStatsLoadStateFunc)).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/values/{dateFrom}/{dateTo}", api.GetMetricValues).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/last-values/{num}", api.GetLastMetricValues).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/last-runs-log/{dateFrom}/{dateTo}", api.GetLastPipelineRunsLog).Methods(http.MethodGet)
	rtr.HandleFunc("/pipelines/{uuid}/logs/{dateFrom}/{dateTo}", api.GetLastPipelineRunsLog).Methods(http.MethodGet)

	rtr.HandleFunc("/private/pipelines/{uuid}/values/avg/{period}/{periodCount}", api.PipelineAvgValueFunc).Methods(http.MethodGet)
	rtr.HandleFunc("/private/pipelines/{uuid}/values/quantile/{level}/{period}/{periodCount}", api.PipelineValueQuantileFunc).Methods(http.MethodGet)
	rtr.HandleFunc("/private/pipelines/{uuid}/duration-stats", api.PipelineDurationStatsQuantileFunc).Methods(http.MethodGet)

	http.Handle("/", rtr)

	server := &http.Server{
		Addr:    config.Cfg().Listen,
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return s.ctx
		},
	}

	go func() {
		<-s.ctx.Done()
		_ = server.Shutdown(s.ctx)
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func NewServer(ctx context.Context, listen string) *Server {
	s := &Server{
		Listen: listen,
		ctx:    ctx,
	}

	return s
}

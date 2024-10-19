package command

import (
	"net"
	"context"
	"net/http"
	"ylem_pipelines/app/envvariable"
	"ylem_pipelines/app/folder"
	"ylem_pipelines/app/task"
	"ylem_pipelines/app/task/result"
	"ylem_pipelines/app/tasktrigger"
	"ylem_pipelines/app/trial"
	"ylem_pipelines/app/pipeline"
	"ylem_pipelines/app/pipelinetemplate"
	"ylem_pipelines/app/dashboard"
	"ylem_pipelines/config"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

var serveHandler cli.ActionFunc = func(c *cli.Context) error {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/pipeline", pipeline.Create).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/trials", trial.CreateTrialOnes).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{uuid}", pipeline.Update).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{uuid}", pipeline.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline/{uuid}/delete", pipeline.Delete).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{uuid}/toggle", pipeline.Toggle).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{uuid}/preview", pipeline.UpdatePreview).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{uuid}/preview", pipeline.FindPreview).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline/{uuid}/clone", pipelinetemplate.ClonePipeline).Methods(http.MethodPost)
	rtr.HandleFunc("/organization/{uuid}/pipelines", pipeline.FindAllInOrganization).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/pipelines/search/{searchString}", pipeline.SearchInOrganization).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/tasks/search/{searchString}", task.SearchInOrganization).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/pipeline-templates/", pipelinetemplate.ListOrganizationTemplates).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/folder/{folderUuid}/pipelines", pipeline.FindAllInOrganizationAndFolder).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/root_folder/pipelines", pipeline.FindAllInOrganizationAndFolder).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/runs_per_month/{type}", pipeline.GetRunsPerOrganizationPerMonth).Methods(http.MethodGet)

	rtr.HandleFunc("/organization/{uuid}/dashboard", dashboard.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/new-grouped-items/{type}/{groupBy}", dashboard.FindNewGroupedItems).Methods(http.MethodGet)

	rtr.HandleFunc("/folder", folder.Create).Methods(http.MethodPost)
	rtr.HandleFunc("/folder/{uuid}", folder.Update).Methods(http.MethodPost)
	rtr.HandleFunc("/folder/{uuid}", folder.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/folder/{folderUuid}/folders", folder.FindAllInOrganizationAndFolder).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/folders", folder.FindAllInOrganizationAndFolder).Methods(http.MethodGet)
	rtr.HandleFunc("/folder/{uuid}/delete", folder.Delete).Methods(http.MethodPost)

	rtr.HandleFunc("/pipeline/{pipelineUuid}/task-trigger", tasktrigger.Create).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task-trigger/{uuid}", tasktrigger.Update).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task-trigger/{uuid}", tasktrigger.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task-trigger/{uuid}/delete", tasktrigger.Delete).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task-triggers", tasktrigger.FindAllInPipeline).Methods(http.MethodGet)

	rtr.HandleFunc("/pipeline/{pipelineUuid}/task", task.Create).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task/{uuid}", task.Update).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task/{uuid}", task.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/task/{uuid}/delete", task.Delete).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/tasks", task.FindAllInPipeline).Methods(http.MethodGet)

	rtr.HandleFunc("/pipeline/{pipelineUuid}/run", result.InitiatePipelineRun).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/run-with-config", result.InitiatePipelineRunWithConfig).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline/{pipelineUuid}/run", result.GetPipelineRunResults).Methods(http.MethodGet)

	rtr.HandleFunc("/envvariable", envvariable.Create).Methods(http.MethodPost)
	rtr.HandleFunc("/envvariable/{uuid}", envvariable.Update).Methods(http.MethodPost)
	rtr.HandleFunc("/envvariable/{uuid}/delete", envvariable.Delete).Methods(http.MethodPost)
	rtr.HandleFunc("/envvariable/{uuid}", envvariable.Find).Methods(http.MethodGet)
	rtr.HandleFunc("/organization/{uuid}/envvariables", envvariable.FindAllInOrganization).Methods(http.MethodGet)

	rtr.HandleFunc("/pipeline-templates/", pipelinetemplate.ListTemplates).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline-templates/", pipelinetemplate.SaveAsTemplate).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline-templates/{templateUuid}/pipelines/", pipelinetemplate.CreateFromTemplate).Methods(http.MethodPost)

	rtr.HandleFunc("/pipeline-templates/{templateUuid}/share-link", pipelinetemplate.GetShareLinkForTemplate).Methods(http.MethodGet)
	rtr.HandleFunc("/pipeline-templates/{templateUuid}/share-link/publish", pipelinetemplate.PublishShareLink).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline-templates/{templateUuid}/share-link/unpublish", pipelinetemplate.UnpublishShareLink).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline-templates/{templateUuid}/share", pipelinetemplate.ShareTemplate).Methods(http.MethodPost)
	rtr.HandleFunc("/pipeline-templates/{templateUuid}/unshare", pipelinetemplate.UnshareTemplate).Methods(http.MethodPost)

	rtr.HandleFunc("/share-links/", pipelinetemplate.ListMyShareLinks).Methods(http.MethodGet)
	rtr.HandleFunc("/share-links/{shareLink}", pipelinetemplate.GetShareLink).Methods(http.MethodGet)

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

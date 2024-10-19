package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services/sql"
	"ylem_integrations/services"

	validation "github.com/go-ozzo/ozzo-validation"
	messaging "github.com/ylem-co/shared-messaging"
	log "github.com/sirupsen/logrus"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

func ShowSQLIntegrationTables(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	uuid := mux.Vars(r)["uuid"]
	err := validation.Validate(uuid, validation.Required, is.UUIDv4)
	if err != nil {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	queryDb := mux.Vars(r)["db"]

	var entity *entities.SQLIntegration
	db := helpers.DbConn()
	defer db.Close()

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionReadList,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if err := entity.Decrypt(r.Context()); err != nil {
		log.Errorf("could not decrypt integration: %s", err.Error())

		return
	}

	password := make([]byte, 0)
	if entity.Password.PlainValue != nil {
		password = entity.Password.PlainValue
	}

	conn, err := sql.CreateSQLIntegrationConnection(
		entity.Type,
		sql.DefaultSQLIntegrationConnectionConfiguration{
			Host:        string(entity.Host.PlainValue),
			Port:        uint16(entity.Port),
			User:        entity.User,
			Password:    string(password),
			Database:    entity.Database,
			SshHost:     string(entity.SshHost.PlainValue),
			SshPort:     uint16(entity.SshPort),
			SshUser:     entity.SshUser,
			ProjectId:   entity.ProjectId,
			Credentials: string(entity.Credentials.PlainValue),
			SslEnabled:  entity.SslEnabled,
			EsVersion:   entity.EsVersion,
		},
	)

	if err != nil {
		log.Errorf("could not collect tables: couldn't create integration connection: %s", err.Error())

		helpers.HttpReturnErrorInternal(w)

		return
	}

	if _, ok := conn.(sql.DescribableConnection); !ok {
		helpers.HttpReturnErrorBadRequest(w, "the integration doesn't support tables describe", &[]string{})

		return
	}

	if entity.ConnectionType == messaging.SQLIntegrationConnectionTypeSsh {
		sshConn, ok := conn.(sql.ViaSshConnection)
		if !ok {
			err = fmt.Errorf("%s connection doesn't support SSH", entity.Type)
		} else {
			err = sshConn.OpenSsh(string(entity.SshHost.PlainValue), uint16(entity.SshPort), entity.SshUser)
		}
	} else {
		err = conn.Open()
	}

	if err != nil {
		log.Errorf("could not collect tables: couldn't open integration connection: %s", err.Error())

		helpers.HttpReturnErrorInternal(w)

		return
	}
	defer conn.Close()

	databases, err := conn.(sql.DescribableConnection).ShowDatabases()
	if err != nil {
		log.Errorf("could not collect tables: couldn't collect databases: %s", err.Error())

		helpers.HttpReturnErrorInternal(w)

		return
	}

	if !CheckSQLDatabaseExists(w, databases, queryDb) {
		return
	}

	tables, err := conn.(sql.DescribableConnection).ShowTables(queryDb)
	if err != nil {
		log.Errorf("could not collect tables: %s", err.Error())

		helpers.HttpReturnErrorInternal(w)

		return
	}

	jsonResponse, err := json.Marshal(map[string]interface{}{
		"items": tables,
	})

	if err != nil {
		log.Errorf("could not marshal tables list: %s", err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

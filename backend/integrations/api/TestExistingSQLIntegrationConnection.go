package api

import (
	"time"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
	"ylem_integrations/services/sql"

	validation "github.com/go-ozzo/ozzo-validation"
	log "github.com/sirupsen/logrus"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

func TestExistingSQLIntegrationConnection(w http.ResponseWriter, r *http.Request) {
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

	var entity *entities.SQLIntegration
	db := helpers.DbConn()
	defer db.Close()

	entity, err = repositories.FindSQLIntegration(db, uuid)

	if err != nil {
		log.Error(err)
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		log.Error("SQL Integration is not found")
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		log.Error("Operation is not permitted")
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if err := entity.Decrypt(r.Context()); err != nil {
		log.Errorf("could not decrypt source: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	password := make([]byte, 0)
	if entity.Password.PlainValue != nil {
		password = entity.Password.PlainValue
	}
	testErr := sql.TestSQLIntegrationConnection(
		entity.Type,
		entity.IsSshConnection(),
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

	if err := entity.Encrypt(r.Context()); err != nil {
		log.Errorf("could not encrypt source: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if testErr != nil {
		helpers.HttpReturnErrorBadRequest(w, testErr.Error(), &[]string{})

		entity.Integration.Status = entities.IntegrationStatusOffline
		err = repositories.UpdateSQLIntegration(db, entity, false)
		if err != nil {
			log.Error(err.Error())
			helpers.HttpReturnErrorInternal(w)
		}

		return
	}

	entity.Integration.Status = entities.IntegrationStatusOnline
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	err = repositories.UpdateSQLIntegration(db, entity, false)
	if err != nil {
		log.Error(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

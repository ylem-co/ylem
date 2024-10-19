package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
)

type HttpSQLIntegrationPrivateResponse struct {
	Uuid             string  `json:"uuid"`
	CreatorUuid      string  `json:"creator_uuid"`
	OrganizationUuid string  `json:"organization_uuid"`
	Status           string  `json:"status"`
	Type             string  `json:"type"`
	Name             string  `json:"name"`
	DataKey          []byte  `json:"data_key"`
	Host             []byte  `json:"host"`
	Port             int     `json:"port,omitempty"`
	User             string  `json:"user,omitempty"`
	Password         []byte  `json:"password"`
	Database         string  `json:"database,omitempty"`
	ConnectionType   string  `json:"connection_type"`
	SslEnabled       bool    `json:"ssl_enabled"`
	SshHost          []byte  `json:"ssh_host,omitempty"`
	SshPort          int     `json:"ssh_port,omitempty"`
	SshUser          string  `json:"ssh_user,omitempty"`
	ProjectId        *string `json:"project_id,omitempty"`
	Credentials      []byte  `json:"credentials,omitempty"`
	EsVersion        *uint8  `json:"es_version,omitempty"`
	UserUpdatedAt    string  `json:"user_updated_at"`
}

func GetSQLIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
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
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	password := make([]byte, 0)
	if entity.Password.EncryptedValue != nil {
		password = entity.Password.EncryptedValue
	}
	response := HttpSQLIntegrationPrivateResponse{
		Uuid:             entity.Integration.Uuid,
		CreatorUuid:      entity.Integration.CreatorUuid,
		OrganizationUuid: entity.Integration.OrganizationUuid,
		Status:           entity.Integration.Status,
		Type:             entity.Type,
		Name:             entity.Integration.Name,
		DataKey:          entity.DataKey.EncryptedValue,
		Host:             entity.Host.EncryptedValue,
		Port:             entity.Port,
		User:             entity.User,
		Password:         password,
		Database:         entity.Database,
		ConnectionType:   entity.ConnectionType,
		SslEnabled:       entity.SslEnabled,
		SshHost:          entity.SshHost.EncryptedValue,
		SshPort:          entity.SshPort,
		SshUser:          entity.SshUser,
		ProjectId:        entity.ProjectId,
		Credentials:      entity.Credentials.EncryptedValue,
		EsVersion:        entity.EsVersion,
		UserUpdatedAt:    entity.Integration.UserUpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(response)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

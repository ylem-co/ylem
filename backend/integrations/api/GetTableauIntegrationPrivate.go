package api

import (
	"encoding/json"
	"net/http"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type GetTableauIntegrationPrivateResponse struct {
	Integration    IntegrationPrivate `json:"integration"`
	Server         string             `json:"server"`
	DataKey        []byte             `json:"data_key"`
	Username       []byte             `json:"username"`
	Password       []byte             `json:"password"`
	Sitename       string             `json:"site_name"`
	ProjectName    string             `json:"project_name"`
	DatasourceName string             `json:"datasource_name"`
	Mode           string             `json:"mode"`
}

func GetTableauIntegrationPrivate(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Tableau
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindTableauIntegration(db, uuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if entity == nil {
		helpers.HttpReturnErrorNotFound(w)

		return
	}

	dataKey, err := decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Username)
	if err != nil {
		return
	}

	_, err = decryptSensitiveData(w, r, entity.Integration.OrganizationUuid, entity.Password)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(GetTableauIntegrationPrivateResponse{
		Integration:    IntegrationPrivate(entity.Integration),
		Server:         entity.Integration.Value,
		Username:       entity.Username.EncryptedValue,
		Password:       entity.Password.EncryptedValue,
		Sitename:       entity.Sitename,
		ProjectName:    entity.ProjectName,
		DatasourceName: entity.DatasourceName,
		Mode:           entity.Mode,
		DataKey:        dataKey,
	})

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

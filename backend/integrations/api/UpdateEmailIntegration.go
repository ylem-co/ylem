package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

type HttpUpdateEmailIntegration struct {
	Name  string `json:"name" valid:"type(string)"`
	Email string `json:"email" valid:"email"`
}

func UpdateEmailIntegration(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	var request HttpUpdateEmailIntegration

	w.Header().Set("Content-Type", "application/json")

	decodeJsonErr := helpers.DecodeJSONBody(w, r, &request)
	if decodeJsonErr != nil {
		rp, _ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		
		_, err := w.Write(rp)
		if err != nil {
			helpers.HttpReturnErrorInternal(w)
			return
		}

		return
	}

	uuid := mux.Vars(r)["uuid"]
	if !govalidator.IsUUIDv4(uuid) {
		helpers.HttpReturnErrorBadUuidRequest(w)

		return
	}

	var entity *entities.Email
	db := helpers.DbConn()
	defer db.Close()

	var err error
	entity, err = repositories.FindEmailIntegration(db, uuid)

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
		services.PermissionActionUpdate,
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	existingEntity, err := repositories.FindIntegrationInOrganizationByValue(db, request.Email, entity.Integration.OrganizationUuid)

	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	if existingEntity != nil && existingEntity.Uuid != uuid {
		helpers.HttpReturnErrorBadRequest(w, "Such Email already exists in the organization", &[]string{})

		return
	}

	isEmailChanged := false

	entity.Integration.Name = request.Name
	if request.Email != entity.Integration.Value {
		isEmailChanged = true
		entity.IsConfirmed = false
	}
	entity.Integration.Value = request.Email
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	if isEmailChanged {
		entity.IsConfirmed = false

		entity.Code = helpers.CreateRandomNumericString(6)
	}

	err = repositories.UpdateEmailIntegration(db, entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	/*if isEmailChanged {
		_, err = services.SendEmailConfirmationEmail(entity.Integration.Value, entity.Integration.Uuid, entity.Code)
		if err != nil {
			fmt.Println(err.Error())
			helpers.HttpReturnErrorInternal(w)

			return
		}
	}*/

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

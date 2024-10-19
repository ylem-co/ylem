package api

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"ylem_integrations/entities"
	"ylem_integrations/helpers"
	"ylem_integrations/repositories"
	"ylem_integrations/services"
)

func ResendConfirmationEmail(w http.ResponseWriter, r *http.Request) {
	user := services.CollectAuthenticationDataByHeader(r.Header.Get("Authorization"))
	if user == nil {
		helpers.HttpReturnErrorUnauthorized(w)

		return
	}

	w.Header().Set("Content-Type", "application/json")

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

	if entity == nil || entity.IsConfirmed {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	canPerformOperation := services.ValidatePermissions(
		user.Uuid,
		entity.Integration.OrganizationUuid,
		services.PermissionActionUpdate, // update?
		services.PermissionResourceTypeIntegration,
		"",
	)
	if !canPerformOperation {
		helpers.HttpReturnErrorForbidden(w)

		return
	}

	if !entity.CanResendEmail() {
		helpers.HttpReturnErrorBadRequest(
			w,
			"Can't resend an email <@todo: when next time can do>",
			nil,
		)

		return
	}

	entity.Code = helpers.CreateRandomNumericString(6)
	entity.RequestedAt = time.Now()
	entity.Integration.UserUpdatedAt = time.Now().Format(helpers.DB_TIME_TIMESTAMP)
	err = repositories.UpdateEmailIntegration(db, entity)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)

		return
	}

	_, err = services.SendEmailConfirmationEmail(entity.Integration.Value, entity.Integration.Uuid, entity.Code)
	if err != nil {
		fmt.Println(err.Error())
		helpers.HttpReturnErrorInternal(w)

		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(entity)
	_, err = w.Write(jsonResponse)
	if err != nil {
		helpers.HttpReturnErrorInternal(w)
		return
	}
}

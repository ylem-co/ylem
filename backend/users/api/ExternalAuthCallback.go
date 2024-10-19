package api

import (
	"encoding/json"
	"fmt"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"ylem_users/services"
	"strconv"
	"strings"
)

func (auth *AuthMiddleware)  ExternalAuthCallback(w http.ResponseWriter, r *http.Request) {
	externalUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Error(err)

		if err.Error() == "could not find a matching session for this request" {
			rp, _ := json.Marshal(map[string]string{"error": "Your request has expired. Please try again."})
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(rp)
			if err != nil {
				log.Error(err)
			}

			return
		}

		helpers.HttpReturnErrorInternal(w)

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	user, ok := repositories.GetUserByEmail(db, externalUser.Email)
	if ok {
		externalAuthCreateJWT(w, user, auth)

		return
	}

	id := externalUser.UserID
	localExternalUser, err := repositories.GetUserByExternalSystemId(db, entities.SourceGoogle, id)
	if err != nil {
		log.Error(err)

		helpers.HttpReturnErrorInternal(w)
		return
	}

	if localExternalUser != nil {
		externalAuthCreateJWT(w, *localExternalUser, auth)

		return
	}

	org, orgEx := repositories.GetOrganization(db)

	srv := services.CreateSignUpSrv(db, r.Context())
	if !orgEx {
		organizationName := strings.Join(strings.Split(externalUser.Email, "@")[1:], "")
		orgExists := repositories.DoesOrganizationExist(db, organizationName, "")

		if orgExists {
			isGenericEmailProvider := services.IsGenericEmailProvider(organizationName)
			if isGenericEmailProvider {
				organizationName = fmt.Sprintf("%s-%s", organizationName, helpers.CreateRandomNumericString(10))
			}

			orgExists = repositories.DoesOrganizationExist(db, organizationName, "")
		}

		if orgExists {
			rp, _ := json.Marshal(map[string]string{"error": "Organization " + organizationName + " already exist. Please contact your organization admin."})
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(rp)
			if err != nil {
				log.Error(err)
			}

			return
		}

		localExternalUser, err = srv.SignUpExternalUser(externalUser, organizationName, entities.SourceGoogle, id, 0)
		if err != nil {
			log.Error(err)

			helpers.HttpReturnErrorInternal(w)
			return
		}

		organization, ok := repositories.GetOrganizationByUserUuid(db, localExternalUser.Uuid)
		if !ok {
			helpers.HttpReturnErrorForbidden(w)
			return
		}

		tokenString := externalAuthCreateJWT(w, *localExternalUser, auth)

		log.Tracef("create default email destination")
		dUuid, ok := services.CreateDefaultEmailDestination(organization.Uuid, localExternalUser.Email, localExternalUser.FirstName, tokenString)

		if ok {
			log.Tracef("create trial db data source")
			sUuid, ok := services.CreateTrialDBDataSource(organization.Uuid, tokenString)
			if ok {
				log.Tracef("create trial pipelines")
				services.TestTrialDBDataSource(sUuid, tokenString)
				_ = services.CreateTrialPipelines(organization.Uuid, dUuid, sUuid, tokenString)
			}
		}
	} else {
		localExternalUser, err = srv.SignUpExternalUser(externalUser, "", entities.SourceGoogle, id, int64(org.Id))
		if err != nil {
			log.Error(err)

			helpers.HttpReturnErrorInternal(w)
			return
		}

		externalAuthCreateJWT(w, *localExternalUser, auth)
	}
}

func externalAuthCreateJWT(w http.ResponseWriter, usr entities.User, auth *AuthMiddleware) string {
	tokenErr, tokenString := CreateJWTToken(usr, auth, nil)
	if tokenErr != nil {
		log.Println(tokenErr)

		return ""
	}

	rp, _ := json.Marshal(map[string]string{"token": tokenString, "uuid": usr.Uuid, "email": usr.Email, "first_name": usr.FirstName, "last_name": usr.LastName, "phone": usr.Phone, "roles": usr.Roles, "is_email_confirmed": strconv.Itoa(usr.IsEmailConfirmed)})
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(rp)
	if err != nil {
		log.Error(err)
	}

	return tokenString
}

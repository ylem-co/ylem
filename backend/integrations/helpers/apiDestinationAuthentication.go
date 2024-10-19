package helpers

import (
	"errors"
	"net/http"
	"ylem_integrations/entities"
)

type AuthTypeCredentialsToProcess struct {
	AuthType      string
	Bearer        *string
	BasicUsername *string
	BasicPassword *string
	HeaderName    *string
	HeaderValue   *string
}

func ProcessHttpAuthTypeRequest(w http.ResponseWriter, Credentials AuthTypeCredentialsToProcess, Entity *entities.Api) error {
	var err error
	if !entities.IsApiAuthTypeSupported(Credentials.AuthType) {
		fields := []string{"auth_type"}
		HttpReturnErrorBadRequest(
			w,
			"Auth type is not supported",
			&fields,
		)

		return errors.New("auth type is not supported")
	}

	switch Credentials.AuthType {
	case entities.ApiAuthTypeNone:
		Entity.SetNoAuth()

		return err
	case entities.ApiAuthTypeBearer:
		err = validateBearer(w, Credentials.Bearer)
		if err == nil {
			Entity.SetBearerAuth(*Credentials.Bearer)
		}

		return err
	case entities.ApiAuthTypeBasic:
		err = validateBasic(w, Credentials.BasicUsername, Credentials.BasicPassword)
		if err == nil {
			Entity.SetBasicAuth(*Credentials.BasicUsername, *Credentials.BasicPassword)
		}

		return err
	case entities.ApiAuthTypeHeader:
		err = validateHeader(w, Credentials.HeaderName, Credentials.HeaderValue)
		if err == nil {
			Entity.SetHeaderAuth(*Credentials.HeaderName, *Credentials.HeaderValue)
		}

		return err
	default:
		panic("should have not reached this stage. check auth_type")
	}
}

func validateHeader(w http.ResponseWriter, name *string, value *string) error {
	if name == nil || len(*name) == 0 || value == nil || len(*value) == 0 {
		fields := []string{"auth_header_name", "auth_header_value"}
		HttpReturnErrorBadRequest(
			w,
			"Header name and/or value should not be empty",
			&fields,
		)

		return errors.New("header name and/or value should not be empty")
	}

	return nil
}

func validateBasic(w http.ResponseWriter, name *string, password *string) error {
	if name == nil || len(*name) == 0 || password == nil || len(*password) == 0 {
		fields := []string{"auth_basic_user_name", "auth_basic_password"}
		HttpReturnErrorBadRequest(
			w,
			"Username and/or password should not be empty",
			&fields,
		)

		return errors.New("username and/or password should not be empty")
	}

	return nil
}

func validateBearer(w http.ResponseWriter, token *string) error {
	if token == nil || len(*token) == 0 {
		fields := []string{"auth_bearer_token"}
		HttpReturnErrorBadRequest(
			w,
			"Bearer token can not be empty",
			&fields,
		)

		return errors.New("bearer token can not be empty")
	}

	return nil
}

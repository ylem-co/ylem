package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/repositories"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	log "github.com/sirupsen/logrus"
)

type HttpUserLogin struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
}

func (auth *AuthMiddleware) LoginUser(w http.ResponseWriter, r *http.Request) {

	var user HttpUserLogin

	w.Header().Set("Content-Type", "application/json")

	err := helpers.DecodeJSONBody(w, r, &user)
	if err != nil {
		rp, _ := json.Marshal(err.Msg)
		w.WriteHeader(err.Status)

		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	errorFields := ValidateHttpUserLogin(user, w)

	if len(errorFields) > 0 {
		rp, _ := json.Marshal(map[string]string{"error": "Invalid fields", "fields": strings.Join(errorFields, ",")})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	db := helpers.DbConn()
	defer db.Close()

	usr, ok := tryToLoginUser(db, user.Email, user.Password)

	var rp []byte
	if !ok {
		rp, _ = json.Marshal(map[string]string{"error": "Invalid User", "fields": ""})
		w.WriteHeader(http.StatusBadRequest)
		
		_, error := w.Write(rp)
		if error != nil {
			log.Error(error)
		}

		return
	}

	vars := r.URL.Query()
	if val, ok := vars["_switch_user"]; ok {
		strRoles := fmt.Sprintf("%v", usr.Roles)
		if strings.Contains(strRoles, entities.ROLE_ALLOWED_TO_SWITCH) {
			email := val[0]
			usr, ok = tryToLoginSwitchUser(db, email)
			if !ok {
				rp, _ = json.Marshal(map[string]string{"error": "Invalid User", "fields": ""})
				w.WriteHeader(http.StatusBadRequest)
				
				_, error := w.Write(rp)
				if error != nil {
					log.Error(error)
				}

				return
			}
		}
	}

	tokenErr, tokenString := CreateJWTToken(usr, auth, nil)
	if tokenErr != nil {
		log.Println(tokenErr)
		return
	}

	rp, _ = json.Marshal(map[string]string{"token": tokenString, "uuid": usr.Uuid, "email": usr.Email, "first_name": usr.FirstName, "last_name": usr.LastName, "phone": usr.Phone, "roles": usr.Roles, "is_email_confirmed": strconv.Itoa(usr.IsEmailConfirmed)})
	w.WriteHeader(http.StatusOK)
	_, error := w.Write(rp)
	if error != nil {
		log.Error(error)
	}
}

func ValidateHttpUserLogin(user HttpUserLogin, w http.ResponseWriter) []string {
	var errorFields []string

	if user.Email == "" || !govalidator.IsEmail(user.Email) {
		errorFields = append(errorFields, "email")
	}

	return errorFields
}

func tryToLoginUser(db *sql.DB, email string, password string) (entities.User, bool) {
	var usr entities.User
	var ok bool

	usr, ok = repositories.GetUserByEmail(db, email)

	if !ok {
		return usr, false
	}

	if bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(password)) != nil {
		return usr, false
	} else {
		return usr, true
	}
}

func tryToLoginSwitchUser(db *sql.DB, email string) (entities.User, bool) {
	var usr entities.User
	var ok bool

	usr, ok = repositories.GetUserByEmail(db, email)

	if !ok {
		return usr, false
	}

	return usr, true
}

func (auth *AuthMiddleware) LogoutUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	rdsUuid := user.RedisUuid
	rdsClient := auth.DbSource.redisDb
	_, err := rdsClient.Del(auth.Ctx, rdsUuid).Result()
	if err == nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

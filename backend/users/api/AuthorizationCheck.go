package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"ylem_users/entities"
	"ylem_users/helpers"
	"ylem_users/services"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type DbSource struct {
	redisDb *redis.Client
}

func (dbSource *DbSource) RdbSource() *redis.Client {
	return dbSource.redisDb
}

func NewDbSource(redisDb *redis.Client) *DbSource {
	return &DbSource{
		redisDb: redisDb,
	}
}

type contextKey string

const authUserKey contextKey = "userKey"

type AuthMiddleware struct {
	//handler  http.Handler
	DbSource *DbSource
	Ctx      context.Context
}

type HttpAuthorizedUser struct {
	Email            string `json:"email"`
	UserUUID         string `json:"uuid"`
	OrganizationUUID string `json:"organization_uuid"`
	DataKey          []byte `json:"data_key"`
}

type CtxAuthorizedUser struct {
	UserId           int
	Email            string
	UserUuid         string
	OrganizationUUID string
	RedisUuid        string
}

func (auth *AuthMiddleware) AuthCheck(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, userId, email, userUuid, orgUuid, rdsUuid := services.InitialAuthorization(w, r)
		if !ok {
			helpers.HttpReturnErrorForbidden(w)
			return
		}

		ctxUser := &CtxAuthorizedUser{UserId: userId, Email: email, UserUuid: userUuid, OrganizationUUID: orgUuid, RedisUuid: rdsUuid}
		rdsClient := auth.DbSource.redisDb
		if rdsClient == nil {
			log.Println("Redis client undefined")
			helpers.HttpReturnErrorForbidden(w)
			return
		}

		_, err := rdsClient.Get(auth.Ctx, rdsUuid).Result()
		if err != nil {
			log.Println("User token in redis is not found. Error: " + err.Error())
			helpers.HttpReturnErrorForbidden(w)
			return
		}
		newRequestCtx := context.WithValue(r.Context(), authUserKey, ctxUser)
		newRequest := r.WithContext(newRequestCtx)
		wrappedHandler.ServeHTTP(w, newRequest)
	})

}

func (auth *AuthMiddleware) IsNotAuthCheck(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, _, _, _, _, _ := services.InitialAuthorization(w, r)
		if ok {
			helpers.HttpReturnErrorForbidden(w)
			return
		}
		wrappedHandler.ServeHTTP(w, r)
	})
}

func AuthorizationCheck(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authUserKey).(*CtxAuthorizedUser)
	userUUID := user.UserUuid
	email := user.Email
	rp, _ := json.Marshal(HttpAuthorizedUser{Email: email, UserUUID: userUUID, OrganizationUUID: user.OrganizationUUID})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(rp)
	if err != nil {
		log.Error(err)
	}
}

func CreateJWTToken(usr entities.User, auth *AuthMiddleware, org *entities.Organization) (error, string) {
	redisUuid := uuid.NewString()
	expTime := time.Now().Add(time.Hour * 24 * 90)
	claims := jwt.MapClaims{
		"uuid":    usr.Uuid,
		"email":   usr.Email,
		"roles":   usr.Roles,
		"exp":     expTime.Unix(),
		"rdsUuid": redisUuid,
	}
	if org != nil {
		claims["organization_uuid"] = org.Uuid
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	pwd, _ := os.Getwd()
	publicKey, errKey := os.ReadFile(pwd + "/config/jwt/private.pem")
	if errKey != nil {
		log.Println(errKey)
		return errKey, ""
	}
	tokenString, _ := token.SignedString(publicKey)

	rdbClient := auth.DbSource.RdbSource()
	now := time.Now()
	rdbErr := rdbClient.Set(auth.Ctx, redisUuid, usr.Uuid, expTime.Sub(now)).Err()
	if rdbErr != nil {
		log.Println(rdbErr)
		return rdbErr, ""
	}

	return nil, tokenString
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")

	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

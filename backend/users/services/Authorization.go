package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"ylem_users/entities"
	"ylem_users/helpers"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type HttpUser struct {
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Password         string `json:"password"`
	ConirmPassword   string `json:"confirm_password"`
	Phone            string `json:"phone"`
	OrganizationName string `json:"organization_name"`
	InvitationKey    string `json:"invitation_key"`
}

func Authorize(r *http.Request) (jwt.MapClaims, bool) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(reqToken) == 0 {
		return nil, false
	}

	tokenString := splitToken[1]
	pwd, _ := os.Getwd()
	signKey, err := os.ReadFile(pwd + "/config/jwt/private.pem")
	if err != nil {
		log.Println(err)
		return nil, false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func InitialAuthorization(w http.ResponseWriter, r *http.Request) (bool, int, string, string, string, string) {
	claims, ok := Authorize(r)
	if !ok {
		return false, 0, "", "", "", ""
	}

	email := fmt.Sprintf("%v", claims["email"])
	var userId = 0

	vars := r.URL.Query()
	if val, ok := vars["_switch_user"]; ok {
		strRoles := fmt.Sprintf("%v", claims["roles"])

		if strings.Contains(strRoles, entities.ROLE_ALLOWED_TO_SWITCH) {
			email = val[0]
		}
	}

	db := helpers.DbConn()
	defer db.Close()

	var userUuid = ""
	orgUuid := ""
	err := db.QueryRow("SELECT u.id, u.uuid, o.uuid AS organization_uuid FROM users u INNER JOIN organizations o ON u.organization_id = o.id WHERE u.email = ?", email).Scan(&userId, &userUuid, &orgUuid)
	if err != nil || userId == 0 {
		log.Println(err, userId)
		return false, 0, "", "", "", ""
	}

	rdsUuid := fmt.Sprintf("%v", claims["rdsUuid"])

	return true, userId, email, userUuid, orgUuid, rdsUuid
}

func OrganizationViewerAuthorization(organizationUuid string, userId int) bool {
	db := helpers.DbConn()
	defer db.Close()

	// validate if it is an organization creator
	var organizationCreatorId = 0
	err := db.QueryRow("SELECT creator_id FROM organizations WHERE uuid = ?", organizationUuid).Scan(&organizationCreatorId)
	if err != nil || organizationCreatorId != userId {
		return false
	}

	return true
}

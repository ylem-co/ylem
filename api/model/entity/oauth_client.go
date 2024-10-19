package entity

import (
	"strconv"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OauthClient struct {
	gorm.Model       `json:"-"`
	Uuid             uuid.UUID `gorm:"unique" json:"uuid"`
	UserUuid         uuid.UUID `json:"user_uuid"`
	OrganizationUuid uuid.UUID `json:"organization_uuid"`
	Name             string    `json:"name"`
	Secret           string    `json:"-"`
	AllowedScopes    string    `json:"-"`
}

// GetID client id
func (c *OauthClient) GetID() string {
	return strconv.FormatInt(int64(c.ID), 10)
}

// GetSecret client secret
func (c *OauthClient) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *OauthClient) GetDomain() string {
	return ""
}

// GetUserID user id
func (c *OauthClient) GetUserID() string {
	return c.Uuid.String()
}

type OauthToken struct {
	gorm.Model
	Uuid                  uuid.UUID `gorm:"unique"`
	OauthClientUuid       uuid.UUID
	OauthClient           *OauthClient `gorm:"foreignKey:OauthClientUuid;references:Uuid"`
	AccessToken           string
	RefreshToken          string
	InternalToken         string
	Scope                 string
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}

// New create to token model instance
func (t *OauthToken) New() oauth2.TokenInfo {
	return NewToken()
}

// GetClientID the client id
func (t *OauthToken) GetClientID() string {
	return t.OauthClientUuid.String()
}

// SetClientID the client id
func (t *OauthToken) SetClientID(clientID string) {
	t.OauthClientUuid = uuid.MustParse(clientID)
}

// GetUserID the user id
func (t *OauthToken) GetUserID() string {
	return t.OauthClient.UserUuid.String()
}

// SetUserID the user id
func (t *OauthToken) SetUserID(userID string) {
}

// GetRedirectURI redirect URI
func (t *OauthToken) GetRedirectURI() string {
	return ""
}

// SetRedirectURI redirect URI
func (t *OauthToken) SetRedirectURI(redirectURI string) {
}

// GetScope get scope of authorization
func (t *OauthToken) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
func (t *OauthToken) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
func (t *OauthToken) GetCode() string {
	return ""
}

// SetCode authorization code
func (t *OauthToken) SetCode(code string) {

}

// GetCodeCreateAt create Time
func (t *OauthToken) GetCodeCreateAt() time.Time {
	return t.CreatedAt
}

// SetCodeCreateAt create Time
func (t *OauthToken) SetCodeCreateAt(createAt time.Time) {
	t.CreatedAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *OauthToken) GetCodeExpiresIn() time.Duration {
	return 0
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *OauthToken) SetCodeExpiresIn(exp time.Duration) {
}

// GetCodeChallenge challenge code
func (t *OauthToken) GetCodeChallenge() string {
	return ""
}

// SetCodeChallenge challenge code
func (t *OauthToken) SetCodeChallenge(code string) {
}

// GetCodeChallengeMethod challenge method
func (t *OauthToken) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod("")
}

// SetCodeChallengeMethod challenge method
func (t *OauthToken) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
}

// GetAccess access Token
func (t *OauthToken) GetAccess() string {
	return t.AccessToken
}

// SetAccess access Token
func (t *OauthToken) SetAccess(access string) {
	t.AccessToken = access
}

// GetAccessCreateAt create Time
func (t *OauthToken) GetAccessCreateAt() time.Time {
	return t.CreatedAt
}

// SetAccessCreateAt create Time
func (t *OauthToken) SetAccessCreateAt(createAt time.Time) {
	t.CreatedAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *OauthToken) GetAccessExpiresIn() time.Duration {
	return time.Duration(t.AccessTokenExpiresIn) * time.Second
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *OauthToken) SetAccessExpiresIn(exp time.Duration) {
	t.AccessTokenExpiresIn = int64(exp.Seconds())
}

func (t *OauthToken) IsAccessExpired() bool {
	return time.Now().After(t.CreatedAt.Add(t.GetAccessExpiresIn()))
}

// GetRefresh refresh Token
func (t *OauthToken) GetRefresh() string {
	return t.RefreshToken
}

// SetRefresh refresh Token
func (t *OauthToken) SetRefresh(refresh string) {
	t.RefreshToken = refresh
}

// GetRefreshCreateAt create Time
func (t *OauthToken) GetRefreshCreateAt() time.Time {
	return t.CreatedAt
}

// SetRefreshCreateAt create Time
func (t *OauthToken) SetRefreshCreateAt(createAt time.Time) {
	t.CreatedAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *OauthToken) GetRefreshExpiresIn() time.Duration {
	return time.Duration(t.RefreshTokenExpiresIn) * time.Second
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *OauthToken) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshTokenExpiresIn = int64(exp.Seconds())
}

func NewToken() *OauthToken {
	return &OauthToken{
		Uuid: uuid.New(),
	}
}

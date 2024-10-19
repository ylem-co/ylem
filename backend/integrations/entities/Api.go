package entities

type Api struct {
	Id                    int64       `json:"-"`
	Integration           Integration `json:"integration"`
	Method                string      `json:"method"`
	AuthType              string      `json:"auth_type"`
	AuthBearerToken       *string     `json:"auth_bearer_token"`
	AuthBasicUserName     *string     `json:"auth_basic_user_name"`
	AuthBasicUserPassword *string     `json:"auth_basic_user_password"`
	AuthHeaderName        *string     `json:"auth_header_name"`
	AuthHeaderValue       *string     `json:"auth_header_value"`
}

const ApiAuthTypeNone = "None"
const ApiAuthTypeBasic = "Basic"
const ApiAuthTypeBearer = "Bearer"
const ApiAuthTypeHeader = "Header"

const IntegrationTypeApi = "api"

func IsApiAuthTypeSupported(AuthType string) bool {
	return map[string]bool{
		ApiAuthTypeNone:   true,
		ApiAuthTypeBasic:  true,
		ApiAuthTypeBearer: true,
		ApiAuthTypeHeader: true,
	}[AuthType]
}

func (a *Api) SetNoAuth() {
	a.AuthType = ApiAuthTypeNone
	a.AuthBearerToken = nil
	a.AuthBasicUserName = nil
	a.AuthBasicUserPassword = nil
	a.AuthHeaderName = nil
	a.AuthHeaderValue = nil
}

func (a *Api) SetBearerAuth(Token string) {
	a.AuthType = ApiAuthTypeBearer
	a.AuthBearerToken = &Token
	a.AuthBasicUserName = nil
	a.AuthBasicUserPassword = nil
	a.AuthHeaderName = nil
	a.AuthHeaderValue = nil
}

func (a *Api) SetBasicAuth(Username string, Password string) {
	a.AuthType = ApiAuthTypeBasic
	a.AuthBearerToken = nil
	a.AuthBasicUserName = &Username
	a.AuthBasicUserPassword = &Password
	a.AuthHeaderName = nil
	a.AuthHeaderValue = nil
}

func (a *Api) SetHeaderAuth(HeaderName string, HeaderValue string) {
	a.AuthType = ApiAuthTypeHeader
	a.AuthBearerToken = nil
	a.AuthBasicUserName = nil
	a.AuthBasicUserPassword = nil
	a.AuthHeaderName = &HeaderName
	a.AuthHeaderValue = &HeaderValue
}

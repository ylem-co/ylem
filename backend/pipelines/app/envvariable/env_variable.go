package envvariable

import "regexp"

type EnvVariable struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	OrganizationUuid string `json:"organization_uuid"`
	Value            string `json:"value"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	IsActive         int8   `json:"-"`
}

func IsEnvVariableValValid(envVariableVal string) bool {
	regExp, _ := regexp.Compile(`^[0-9_-]*[0-9.]*[a-zA-Z]*[a-zA-Z0-9_-]*$`)
	return regExp.MatchString(envVariableVal)
}

func IsEnvVariableNameValid(envVariableName string) bool {
	regExp, _ := regexp.Compile(`^[0-9_-]*[a-zA-Z]*[a-zA-Z0-9_-]*$`)
	return regExp.MatchString(envVariableName)
}

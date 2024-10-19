package entities

type Salesforce struct {
	Id                      int64                   `json:"-"`
	Integration             Integration             `json:"integration"`
	SalesforceAuthorization SalesforceAuthorization `json:"authorization"`
}

const IntegrationTypeSalesforce = "salesforce"

package entities

type Jira struct {
	Id                int64             `json:"-"`
	Integration       Integration       `json:"integration"`
	JiraAuthorization JiraAuthorization `json:"authorization"`
	IssueType         string            `json:"issue_type"`
}

const IntegrationTypeJira = "jira"

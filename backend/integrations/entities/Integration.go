package entities

type Integration struct {
	Id               int64   `json:"-"`
	Uuid             string  `json:"uuid"`
	CreatorUuid      string  `json:"-"`
	OrganizationUuid string  `json:"-"`
	Status           string  `json:"status"`
	Type             string  `json:"type"`
	IoType           string  `json:"io_type"`
	Name             string  `json:"name"`
	Value            string  `json:"value"`
	UserUpdatedAt    string  `json:"user_updated_at"`
}

const IntegrationStatusNew = "new"
const IntegrationStatusOnline = "online"
const IntegrationStatusOffline = "offline"

const IntegrationIoTypeAll = "all"
const IntegrationIoTypeSQL = "sql"
const IntegrationIoTypeRead = "read"
const IntegrationIoTypeWrite = "write"
const IntegrationIoTypeReadWrite = "read-write"

func IsIntegrationStatusSupported(Status string) bool {
	return Status == IntegrationStatusOnline || Status == IntegrationStatusOffline
}

func IsIoTypeSupported(Type string) bool {
	return map[string]bool{
		IntegrationIoTypeAll:        true,
		IntegrationIoTypeSQL:        true,
		IntegrationIoTypeRead:       true,
		IntegrationIoTypeWrite:      true,
		IntegrationIoTypeReadWrite:  true,
	}[Type]
}
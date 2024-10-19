package task

import (
	"regexp"
)

type Task struct {
	Id               int64       `json:"-"`
	Uuid             string      `json:"uuid"`
	Name             string      `json:"name"`
	Severity         string      `json:"severity"`
	PipelineUuid     string      `json:"pipeline_uuid"`
	OrganizationUuid string      `json:"-"`
	PipelineId       int64       `json:"-"`
	Type             string      `json:"type"`
	ImplementationId int64       `json:"-"`
	Implementation   interface{} `json:"implementation"`
	IsActive         int8        `json:"-"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
}

type SearchedTask struct {
	Id               int64       `json:"-"`
	Uuid             string      `json:"uuid"`
	Name             string      `json:"name"`
	PipelineUuid     string      `json:"pipeline_uuid"`
	OrganizationUuid string      `json:"-"`
	PipelineId       int64       `json:"-"`
	Type             string      `json:"type"`
	FolderUuid       string      `json:"folder_uuid"`
	FolderId         int64       `json:"-"`
	IsActive         int8        `json:"-"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
}

type Condition struct {
	Id         int64  `json:"-"`
	Uuid       string `json:"uuid"`
	Expression string `json:"expression"`
	IsActive   int8   `json:"-"`
	CreatedAt  string `json:"-"`
	UpdatedAt  string `json:"-"`
}

type ForEach struct {
	Id        int64  `json:"-"`
	Uuid      string `json:"uuid"`
	IsActive  int8   `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

type Aggregator struct {
	Id           int64  `json:"-"`
	Uuid         string `json:"uuid"`
	Expression   string `json:"expression"`
	VariableName string `json:"variable_name"`
	IsActive     int8   `json:"-"`
	CreatedAt    string `json:"-"`
	UpdatedAt    string `json:"-"`
}

type Processor struct {
	Id           int64  `json:"-"`
	Uuid         string `json:"uuid"`
	Expression   string `json:"expression"`
	Strategy     string `json:"strategy"`
	IsActive     int8   `json:"-"`
	CreatedAt    string `json:"-"`
	UpdatedAt    string `json:"-"`
}

type Query struct {
	Id         int64  `json:"-"`
	Uuid       string `json:"uuid"`
	SQLQuery   string `json:"sql_query"`
	SourceUuid string `json:"source_uuid"`
	IsActive   int8   `json:"-"`
	CreatedAt  string `json:"-"`
	UpdatedAt  string `json:"-"`
}

type Transformer struct {
	Id                  int64  `json:"-"`
	Uuid                string `json:"uuid"`
	Type                string `json:"type"`
	JsonQueryExpression string `json:"json_query_expression"`
	Delimiter           string `json:"delimiter"`
	CastToType          string `json:"cast_to_type"`
	DecodeFormat        string `json:"decode_format"`
	EncodeFormat        string `json:"encode_format"`
	IsActive            int8   `json:"-"`
	CreatedAt           string `json:"-"`
	UpdatedAt           string `json:"-"`
}

type Notification struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Type             string `json:"type"`
	Body             string `json:"body"`
	AttachedFileName string `json:"attached_file_name"`
	DestinationUuid  string `json:"destination_uuid"`
	IsActive         int8   `json:"-"`
	CreatedAt        string `json:"-"`
	UpdatedAt        string `json:"-"`
}

type ApiCall struct {
	Id               int64  `json:"-"`
	Uuid             string `json:"uuid"`
	Type             string `json:"type"`
	Payload          string `json:"payload"`
	QueryString      string `json:"query_string"`
	Headers          string `json:"headers"`
	AttachedFileName string `json:"attached_file_name"`
	DestinationUuid  string `json:"destination_uuid"`
	IsActive         int8   `json:"-"`
	CreatedAt        string `json:"-"`
	UpdatedAt        string `json:"-"`
}

type Merge struct {
	Id         int64  `json:"-"`
	Uuid       string `json:"uuid"`
	FieldNames string `json:"field_names"`
	IsActive   int8   `json:"-"`
	CreatedAt  string `json:"-"`
	UpdatedAt  string `json:"-"`
}

type Filter struct {
	Id         int64  `json:"-"`
	Uuid       string `json:"uuid"`
	Expression string `json:"expression"`
	IsActive   int8   `json:"-"`
	CreatedAt  string `json:"-"`
	UpdatedAt  string `json:"-"`
}

type RunPipeline struct {
	Id           int64  `json:"-"`
	Uuid         string `json:"uuid"`
	PipelineUuid string `json:"pipeline_uuid"`
	IsActive     int8   `json:"-"`
	CreatedAt    string `json:"-"`
	UpdatedAt    string `json:"-"`
}

type ExternalTrigger struct {
	Id        int64  `json:"-"`
	Uuid      string `json:"uuid"`
	TestData  string `json:"test_data"`
	IsActive  int8   `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

type Code struct {
	Id        int64  `json:"-"`
	Uuid      string `json:"uuid"`
	Code      string `json:"code"`
	Type      string `json:"type"`
	IsActive  int8   `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

type Gpt struct {
	Id        int64  `json:"-"`
	Uuid      string `json:"uuid"`
	Prompt    string `json:"prompt"`
	IsActive  int8   `json:"-"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

const TaskIsActive = 1
const TaskIsNotActive = 0

const TaskTypeQuery = "query"
const TaskTypeCondition = "condition"
const TaskTypeAggregator = "aggregator"
const TaskTypeProcessor = "processor"
const TaskTypeTransformer = "transformer"
const TaskTypeNotification = "notification"
const TaskTypeApiCall = "api_call"
const TaskTypeForEach = "for_each"
const TaskTypeMerge = "merge"
const TaskTypeFilter = "filter"
const TaskTypeExternalTrigger = "external_trigger"
const TaskTypeCode = "code"
const TaskTypeGpt = "gpt"
const TaskTypeRunPipeline = "run_pipeline"

const TaskSeverityLowest = "lowest"
const TaskSeverityLow = "low"
const TaskSeverityMedium = "medium"
const TaskSeverityHigh = "high"
const TaskSeverityCritical = "critical"

const NotificationTypeSms = "sms"
const NotificationTypeEmail = "email"
const NotificationTypeSlack = "slack"
const NotificationTypeJira = "jira"
const NotificationTypeIncidentIo = "incidentio"
const NotificationTypeTableau = "tableau"
const NotificationTypeHubspot = "hubspot"
const NotificationTypeGoogleSheets = "google-sheets"
const NotificationTypeSalesforce = "salesforce"
const NotificationTypeOpsgenie = "opsgenie"
const NotificationTypeJenkins = "jenkins"

const ApiCallTypeApi = "api"

const TransformerTypeStrSplit = "str_split"
const TransformerTypeExtractFromJSON = "extract_from_json"
const TransformerTypeCastTo = "cast_to"
const TransformerTypeEncode = "encode_to"

const TransformerTypeCastToString = "string"
const TransformerTypeCastToInteger = "integer"

const TransformerTypeEncodeToXML = "XML"
const TransformerTypeEncodeToCSV = "CSV"

const ProcessorStrategyInclusive = "inclusive"
const ProcessorStrategyExclusive = "exclusive"

func IsTypeSupported(Type string) bool {
	return map[string]bool{
		TaskTypeQuery:           true,
		TaskTypeCondition:       true,
		TaskTypeAggregator:      true,
		TaskTypeTransformer:     true,
		TaskTypeNotification:    true,
		TaskTypeApiCall:         true,
		TaskTypeForEach:         true,
		TaskTypeMerge:           true,
		TaskTypeFilter:          true,
		TaskTypeExternalTrigger: true,
		TaskTypeCode:            true,
		TaskTypeGpt:             true,
		TaskTypeRunPipeline:     true,
		TaskTypeProcessor:       true,
	}[Type]
}

func IsSeveritySupported(Severity string) bool {
	return true
}

func IsNotificationTypeSupported(Type string) bool {
	return map[string]bool{
		NotificationTypeSms:          true,
		NotificationTypeEmail:        true,
		NotificationTypeSlack:        true,
		NotificationTypeJira:         true,
		NotificationTypeIncidentIo:   true,
		NotificationTypeTableau:      true,
		NotificationTypeHubspot:      true,
		NotificationTypeSalesforce:   true,
		NotificationTypeOpsgenie:     true,
		NotificationTypeJenkins:      true,
		NotificationTypeGoogleSheets: true,
	}[Type]
}

func IsProcessorStrategySupported(Strategy string) bool {
	return map[string]bool{
		ProcessorStrategyExclusive: true,
		ProcessorStrategyInclusive: true,
	}[Strategy]
}

func IsApiCallTypeSupported(Type string) bool {
	return map[string]bool{
		ApiCallTypeApi: true,
	}[Type]
}

func IsTransformerTypeSupported(Type string) bool {
	return map[string]bool{
		TransformerTypeStrSplit:        true,
		TransformerTypeExtractFromJSON: true,
		TransformerTypeCastTo:          true,
		TransformerTypeEncode:          true,
	}[Type]
}

func IsTransformerValid(Transformer *HttpApiNewTransformer) bool {
	if Transformer.Type == TransformerTypeCastTo {
		return map[string]bool{
			TransformerTypeCastToString:  true,
			TransformerTypeCastToInteger: true,
		}[Transformer.CastToType]
	}

	if Transformer.Type == TransformerTypeEncode {
		return map[string]bool{
			TransformerTypeEncodeToXML: true,
			TransformerTypeEncodeToCSV: true,
		}[Transformer.EncodeFormat]
	}

	return true
}

func IsUpdatedTransformerValid(Transformer *HttpApiUpdatedTransformer) bool {
	if Transformer.Type == TransformerTypeCastTo {
		return map[string]bool{
			TransformerTypeCastToString:  true,
			TransformerTypeCastToInteger: true,
		}[Transformer.CastToType]
	}

	if Transformer.Type == TransformerTypeEncode {
		return map[string]bool{
			TransformerTypeEncodeToXML: true,
			TransformerTypeEncodeToCSV: true,
		}[Transformer.EncodeFormat]
	}

	return true
}

func IsQuerySafe(query string) bool {
	queryRegex := regexp.MustCompile(`(?i)(DROP|TRUNCATE|DELETE)\s`)
	return !queryRegex.MatchString(query)
}

package task

type HttpApiNewTask struct {
	Name         string                  `json:"name" valid:"type(string)"`
	Type         string                  `json:"type" valid:"type(string)"`
	Condition    *HttpApiNewCondition    `json:"condition" valid:"optional"`
	Aggregator   *HttpApiNewAggregator   `json:"aggregator" valid:"optional"`
	Transformer  *HttpApiNewTransformer  `json:"transformer" valid:"optional"`
	Query        *HttpApiNewQuery        `json:"query" valid:"optional"`
	Notification *HttpApiNewNotification `json:"notification" valid:"optional"`
	ApiCall      *HttpApiNewApiCall      `json:"api_call" valid:"optional"`
	ForEach      *HttpApiNewForEach      `json:"for_each" valid:"optional"`
	Merge        *HttpApiNewMerge        `json:"merge" valid:"optional"`
	Filter       *HttpApiNewFilter       `json:"filter" valid:"optional"`
	RunPipeline  *HttpApiNewRunPipeline  `json:"run_pipeline" valid:"optional"`
	Code         *HttpApiNewCode         `json:"code" valid:"optional"`
	Gpt          *HttpApiNewGpt          `json:"gpt" valid:"optional"`
	Processor    *HttpApiNewProcessor    `json:"processor" valid:"optional"`
}

type HttpApiUpdatedTask struct {
	Name            string                         `json:"name" valid:"type(string)"`
	Severity        string                         `json:"severity" valid:"type(string)"`
	Condition       *HttpApiUpdatedCondition       `json:"condition" valid:"optional"`
	Aggregator      *HttpApiUpdatedAggregator      `json:"aggregator" valid:"optional"`
	Transformer     *HttpApiUpdatedTransformer     `json:"transformer" valid:"optional"`
	Query           *HttpApiUpdatedQuery           `json:"query" valid:"optional"`
	Notification    *HttpApiUpdatedNotification    `json:"notification" valid:"optional"`
	ApiCall         *HttpApiUpdatedApiCall         `json:"api_call" valid:"optional"`
	ForEach         *HttpApiUpdatedForEach         `json:"for_each" valid:"optional"`
	Merge           *HttpApiUpdatedMerge           `json:"merge" valid:"optional"`
	Filter          *HttpApiUpdatedFilter          `json:"filter" valid:"optional"`
	RunPipeline     *HttpApiUpdatedRunPipeline     `json:"run_pipeline" valid:"optional"`
	ExternalTrigger *HttpApiUpdatedExternalTrigger `json:"external_trigger" valid:"optional"`
	Code            *HttpApiUpdatedCode            `json:"code" valid:"optional"`
	Gpt             *HttpApiUpdatedGpt             `json:"gpt" valid:"optional"`
	Processor       *HttpApiUpdatedProcessor       `json:"processor" valid:"optional"`
}

type HttpApiNewForEach struct{}

type HttpApiUpdatedForEach struct{}

type HttpApiNewCondition struct {
	Expression string `json:"expression" valid:"type(string)"`
}

type HttpApiUpdatedCondition struct {
	Expression string `json:"expression" valid:"type(string)"`
}

type HttpApiNewAggregator struct {
	Expression   string `json:"expression" valid:"type(string)"`
	VariableName string `json:"variable_name" valid:"type(string), optional"`
}

type HttpApiUpdatedAggregator struct {
	Expression   string `json:"expression" valid:"type(string)"`
	VariableName string `json:"variable_name" valid:"type(string), optional"`
}

type HttpApiNewProcessor struct {
	Expression   string `json:"expression" valid:"type(string)"`
	Strategy     string `json:"strategy" valid:"type(string)"`
}

type HttpApiUpdatedProcessor struct {
	Expression   string `json:"expression" valid:"type(string)"`
	Strategy     string `json:"strategy" valid:"type(string)"`
}

type HttpApiNewTransformer struct {
	Type                string `json:"type" valid:"type(string)"`
	JsonQueryExpression string `json:"json_query_expression" valid:"type(string), optional"`
	Delimiter           string `json:"delimiter" valid:"type(string), optional"`
	CastToType          string `json:"cast_to_type" valid:"type(string), optional"`
	DecodeFormat        string `json:"decode_format" valid:"type(string), optional"`
	EncodeFormat        string `json:"encode_format" valid:"type(string), optional"`
}

type HttpApiUpdatedTransformer struct {
	Type                string `json:"type" valid:"type(string)"`
	JsonQueryExpression string `json:"json_query_expression" valid:"type(string), optional"`
	Delimiter           string `json:"delimiter" valid:"type(string), optional"`
	CastToType          string `json:"cast_to_type" valid:"type(string), optional"`
	DecodeFormat        string `json:"decode_format" valid:"type(string), optional"`
	EncodeFormat        string `json:"encode_format" valid:"type(string), optional"`
}

type HttpApiNewNotification struct {
	Type             string `json:"type" valid:"type(string)"`
	Body             string `json:"body" valid:"type(string), optional"`
	AttachedFileName string `json:"attached_file_name" valid:"type(string), optional"`
	DestinationUuid  string `json:"destination_uuid" valid:"uuidv4"`
}

type HttpApiUpdatedNotification struct {
	Type             string `json:"type" valid:"type(string)"`
	Body             string `json:"body" valid:"type(string), optional"`
	AttachedFileName string `json:"attached_file_name" valid:"type(string), optional"`
	DestinationUuid  string `json:"destination_uuid" valid:"uuidv4"`
}

type HttpApiNewQuery struct {
	SQLQuery   string `json:"sql_query" valid:"type(string)"`
	SourceUuid string `json:"source_uuid" valid:"uuidv4"`
}

type HttpApiUpdatedQuery struct {
	SQLQuery   string `json:"sql_query" valid:"type(string)"`
	SourceUuid string `json:"source_uuid" valid:"uuidv4"`
}

type HttpApiNewApiCall struct {
	Type             string `json:"type" valid:"type(string)"`
	Payload          string `json:"payload" valid:"type(string), optional"`
	QueryString      string `json:"query_string" valid:"type(string), optional"`
	Headers          string `json:"headers" valid:"type(string), optional"`
	AttachedFileName string `json:"attached_file_name" valid:"type(string), optional"`
	DestinationUuid  string `json:"destination_uuid" valid:"uuidv4"`
}

type HttpApiUpdatedApiCall struct {
	Type             string `json:"type" valid:"type(string)"`
	Payload          string `json:"payload" valid:"type(string), optional"`
	QueryString      string `json:"query_string" valid:"type(string), optional"`
	Headers          string `json:"headers" valid:"type(string), optional"`
	AttachedFileName string `json:"attached_file_name" valid:"type(string), optional"`
	DestinationUuid  string `json:"destination_uuid" valid:"uuidv4"`
}

type HttpApiNewMerge struct {
	FieldNames string `json:"field_names" valid:"type(string),optional"`
}

type HttpApiUpdatedMerge struct {
	FieldNames string `json:"field_names" valid:"type(string),optional"`
}

type HttpApiNewFilter struct {
	Expression string `json:"expression" valid:"type(string),optional"`
}

type HttpApiNewCode struct {
	Code string `json:"code" valid:"type(string),optional"`
	Type string `json:"type" valid:"type(string),optional"`
}

type HttpApiNewGpt struct {
	Prompt string `json:"prompt" valid:"type(string)"`
}

type HttpApiUpdatedGpt struct {
	Prompt string `json:"prompt" valid:"type(string),optional"`
}

type HttpApiUpdatedFilter struct {
	Expression string `json:"expression" valid:"type(string),optional"`
}

type HttpApiUpdatedCode struct {
	Code string `json:"code" valid:"type(string),optional"`
	Type string `json:"type" valid:"type(string),optional"`
}

type HttpApiNewRunPipeline struct {
	PipelineUuid string `json:"pipeline_uuid" valid:"type(string)"`
}

type HttpApiUpdatedRunPipeline struct {
	PipelineUuid string `json:"pipeline_uuid" valid:"type(string)"`
}

type HttpApiUpdatedExternalTrigger struct {
	TestData string `json:"test_data" valid:"type(string),optional"`
}

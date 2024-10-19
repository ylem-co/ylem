package common

const (
	PipelineTypeGeneric = "generic"
	PipelineTypeMetric  = "metric"
)

func IsTypeSupported(Type string) bool {
	return map[string]bool{
		PipelineTypeGeneric:      true,
		PipelineTypeMetric:       true,
	}[Type]
}

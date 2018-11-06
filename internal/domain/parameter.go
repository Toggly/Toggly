package domain

// ParameterType type
type ParameterType string

// Parameter types enum
const (
	ParameterBool   ParameterType = "bool"
	ParameterString ParameterType = "string"
	ParameterInt    ParameterType = "int"
	ParameterEnum   ParameterType = "enum"
)

// ParameterCode type
type ParameterCode string

// Parameter represents a flag data structure
type Parameter struct {
	Code        ParameterCode `json:"code"`
	Description string        `json:"description"`
	Type        ParameterType `json:"type"`
	Value       interface{}
}

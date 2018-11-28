package domain

// ParameterType type
type ParameterType string

// Parameter types enum
const (
	ParameterBool   ParameterType = "bool"
	ParameterString ParameterType = "string"
	ParameterInt    ParameterType = "int"
)

// ParameterCode type
type ParameterCode string

// Parameter represents a flag data structure
type Parameter struct {
	Code          ParameterCode `json:"code"`
	Description   string        `json:"description"`
	Type          ParameterType `json:"type"`
	Value         interface{}   `json:"value"`
	AllowedValues []interface{} `json:"allowed_values,omitempty" bson:"allowed_values,omitempty"`
}

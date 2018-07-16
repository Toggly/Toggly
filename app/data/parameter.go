package data

// ParameterType type
type ParameterType int

// Parameter types enum
const (
	ParameterBool ParameterType = iota
	ParameterString
	ParameterInt
	ParameterEnum
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

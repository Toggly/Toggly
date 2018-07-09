package data

// Object describes object configuration
type Object struct {
	Code        CodeType    `json:"code"`
	Description string      `json:"description"`
	Inherits    CodeType    `json:"inherits"`
	Parameters  []Parameter `json:"parameters"`
}

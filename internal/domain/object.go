package domain

// ObjectCode type
type ObjectCode string

// Object describes object configuration
type Object struct {
	Code        ObjectCode  `json:"code"`
	Description string      `json:"description"`
	Inherits    ObjectCode  `json:"inherits"`
	Parameters  []Parameter `json:"parameters"`
}

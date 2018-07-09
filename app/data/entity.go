package data

import "time"

// ObjectID is an unique object identifier
type ObjectID string

// Environment represents an environment data structure
type Environment struct {
	ID        ObjectID  `json:"_id"`
	Name      string    `json:"name"`
	Protected bool      `json:"protected"`
	RegDate   time.Time `json:"reg_date,omitempty"`
}

// OptionType type
type OptionType int

// Option types enum
const (
	OptionBool OptionType = iota
	OptionString
	OptionInt
	OptionEnum
)

// ParameterCode type
type ParameterCode string

// Parameter represents a flag data structure
type Parameter struct {
	Code        ParameterCode `json:"code"`
	Description string        `json:"description"`
	Type        OptionType    `json:"type"`
	Value       interface{}
}

// ObjectCode type
type ObjectCode string

// Object describes object configuration
type Object struct {
	Code       ObjectCode
	Name       string
	Overrides  ObjectCode
	Parameters []Parameter
}

// Account represents an account data structure
type Account struct {
	ID      ObjectID  `json:"id"`
	Name    string    `json:"name"`
	OAuthID string    `json:"oauth_id"`
	RegDate time.Time `json:"reg_date"`
}

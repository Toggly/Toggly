package rest

import "time"

// ObjectID is an unique object identifier
type ObjectID string

// Project represents a project data structure
type Project struct {
	ID             ObjectID  `json:"id"`
	Name           string    `json:"name,omitempty"`
	Collaboratives []int     `json:"collaboratives,omitempty"`
	RegDate        time.Time `json:"reg_date,omitempty"`
	Status         int       `json:"status"`
}

// Environment represents an environment data structure
type Environment struct {
	ID        ObjectID `json:"id"`
	Name      string   `json:"name"`
	Protected bool     `json:"protected"`
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

// Object describes dictionary object
type Object struct {
	Code      ObjectCode
	Name      string
	Overrides ObjectCode
	Props     map[string]interface{}
}

// Account represents an account data structure
type Account struct {
	ID      ObjectID  `json:"id"`
	Name    string    `json:"name"`
	OAuthID string    `json:"oauth_id"`
	RegDate time.Time `json:"reg_date"`
}

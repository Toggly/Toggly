package rest

import "time"

// ObjectID is an unique object identifier
type ObjectID string

// Account represents an account data structure
type Account struct {
	ID      ObjectID  `json:"id"`
	Name    string    `json:"name"`
	OAuthID string    `json:"oauth_id"`
	RegDate time.Time `json:"reg_date"`
}

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

// FlagType type
type FlagType int

// flag types enum
const (
	FlagBool FlagType = iota
	FlagString
	FlagInt
	FlagEnum
)

// Flag represents a flag data structure
type Flag struct {
	ID        ObjectID `json:"id"`
	Name      string   `json:"name"`
	Type      FlagType `json:"type"`
	Protected bool     `json:"protected"`
}

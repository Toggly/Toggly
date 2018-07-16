package data

import "time"

// EnvironmentCode type
type EnvironmentCode string

// Environment represents an environment data structure
type Environment struct {
	Code        EnvironmentCode `json:"code"`
	Description string          `json:"description"`
	Protected   bool            `json:"protected"`
	RegDate     time.Time       `json:"reg_date,omitempty"`
}

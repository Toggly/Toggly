package data

import "time"

// Environment represents an environment data structure
type Environment struct {
	Code        CodeType  `json:"code"`
	Description string    `json:"description"`
	Protected   bool      `json:"protected"`
	RegDate     time.Time `json:"reg_date,omitempty"`
}

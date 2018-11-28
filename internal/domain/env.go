package domain

import "time"

// EnvironmentCode type
type EnvironmentCode string

// Environment represents an environment data structure
type Environment struct {
	OwnerID     string          `json:"owner" bson:"owner"`
	ProjectCode ProjectCode     `json:"project_code" bson:"project_code"`
	Code        EnvironmentCode `json:"code"`
	Description string          `json:"description"`
	Protected   bool            `json:"protected"`
	RegDate     time.Time       `json:"reg_date" bson:"reg_date"`
	Inherits    EnvironmentCode `json:"inherits"`
}

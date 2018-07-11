package data

import "time"

// Project represents a project data structure
type Project struct {
	ID          string    `json:"id" bson:"_id"`
	Code        CodeType  `json:"code"`
	Description string    `json:"description,omitempty"`
	RegDate     time.Time `json:"reg_date,omitempty"`
	Status      int       `json:"status"`
}

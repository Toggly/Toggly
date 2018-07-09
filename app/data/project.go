package data

import "time"

// ProjectCode type
type ProjectCode string

// Project represents a project data structure
type Project struct {
	Code        ProjectCode `json:"code"`
	Description string      `json:"description,omitempty"`
	RegDate     time.Time   `json:"reg_date,omitempty"`
	Status      int         `json:"status"`
}

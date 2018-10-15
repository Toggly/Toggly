package domain

import (
	"time"
)

// ProjectCode type
type ProjectCode string

// ProjectStatus type
type ProjectStatus string

// ProjectStatus enum
const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusDisabled ProjectStatus = "disabled"
)

// Project represents a project data structure
type Project struct {
	OwnerID     string        `json:"owner" bson:"owner"`
	Code        ProjectCode   `json:"code" bson:"code"`
	Description string        `json:"description,omitempty"`
	RegDate     time.Time     `json:"reg_date,omitempty"`
	Status      ProjectStatus `json:"status"`
}

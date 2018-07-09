package data

import "time"

// Account represents an account data structure
type Account struct {
	ID      ObjectID  `json:"id"`
	Name    string    `json:"name"`
	OAuthID string    `json:"oauth_id"`
	RegDate time.Time `json:"reg_date"`
}

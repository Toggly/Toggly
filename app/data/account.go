package data

import "time"

// Account represents an account data structure
type Account struct {
	ID      ObjectID  `json:"id" bson:"_id"`
	Name    string    `json:"name"`
	RegDate time.Time `json:"reg_date"`
}

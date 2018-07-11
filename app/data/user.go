package data

// User represents user data
type User struct {
	ID      ObjectID `json:"id" bson:"_id"`
	Account Account  `json:"account"`
	Name    string   `json:"name"`
}

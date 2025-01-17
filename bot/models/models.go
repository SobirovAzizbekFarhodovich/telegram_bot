package models

type Password struct {
	UserID string `json:"user_id"`
	Site   string `json:"site"`
	Password string `json:"password"`
}



type Passwordswagger struct {
	ID string `json:"id" db:"id"`
	Site string `json:"site" db:"site"`
	Password string `json:"password" db:"password"`
}
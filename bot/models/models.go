package models

type Password struct {
    ID       string `json:"id" db:"id"`
    Site     string `json:"site" db:"site"`
    Password string `json:"password" db:"password"`
    UserID   string `json:"user_id" db:"user_id"` 
}


type Passwordswagger struct {
	ID string `json:"id" db:"id"`
	Site string `json:"site" db:"site"`
	Password string `json:"password" db:"password"`
}
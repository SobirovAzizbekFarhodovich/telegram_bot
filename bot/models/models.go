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

type APIResponse struct {
	Timestamp   string      `json:"timestamp" example:"2025-01-17T20:58:06Z"`
	RequestURL  string      `json:"request_url" example:"/password/get_password"`
	Message     string      `json:"message" example:"Passwords retrieved successfully"`
	Reason      string      `json:"reason,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

package models

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password,omitempty"`
}

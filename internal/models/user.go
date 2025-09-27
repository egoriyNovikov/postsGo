package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (u *User) CreateUser(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (username, email, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5)", u.Username, u.Email, u.CreatedAt, u.UpdatedAt, u.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

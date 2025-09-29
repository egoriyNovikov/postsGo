package models

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID         int       `json:"id"`
	User_id    int       `json:"user_id"`
	Token      string    `json:"token"`
	Expires_at time.Time `json:"expires_at"`
}

func (a *Auth) CreateAuth(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO auth (user_id, token, expires_at) VALUES ($1, $2, $3)", a.User_id, a.Token, a.Expires_at)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) GetAuth(db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM auth WHERE id = $1", a.ID)
	if err := row.Scan(&a.ID, &a.User_id, &a.Token, &a.Expires_at); err != nil {
		return err
	}
	return nil
}

func (a *Auth) GetAuthByToken(db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM auth WHERE token = $1", a.Token)
	if err := row.Scan(&a.ID, &a.User_id, &a.Token, &a.Expires_at); err != nil {
		return err
	}
	return nil
}

func (a *Auth) DeleteAuth(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM auth WHERE token = $1", a.Token)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) DeleteExpiredTokens(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM auth WHERE expires_at < NOW()")
	if err != nil {
		return err
	}
	return nil
}

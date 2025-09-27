package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (p *Post) CreatePost(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO posts (title, content, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5)", p.Title, p.Content, p.CreatedAt, p.UpdatedAt, p.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetPost(db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM posts WHERE id = $1", p.ID)
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
		return err
	}
	return nil
}

func (p *Post) GetAllPosts(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM posts WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *Post) UpdatePost(db *sql.DB) error {
	_, err := db.Exec("UPDATE posts SET title = $1, content = $2, updated_at = $3 WHERE id = $4", p.Title, p.Content, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = $1", p.ID)
	if err != nil {
		return err
	}
	return nil
}

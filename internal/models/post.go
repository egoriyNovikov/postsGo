package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (p *Post) CreatePost(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO posts (user_id, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", p.UserID, p.Title, p.Content, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetPost(db *sql.DB) error {
	row := db.QueryRow("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM posts WHERE id = $1 AND deleted_at IS NULL", p.ID)
	var deletedAt sql.NullTime
	err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &deletedAt)
	if err != nil {
		return err
	}
	if deletedAt.Valid {
		p.DeletedAt = &deletedAt.Time
	}
	return nil
}

func (p *Post) GetAllPosts(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM posts WHERE deleted_at IS NULL ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *Post) UpdatePost(db *sql.DB) error {
	_, err := db.Exec("UPDATE posts SET title = $1, content = $2, updated_at = $3 WHERE id = $4 AND user_id = $5 AND deleted_at IS NULL", p.Title, p.Content, p.UpdatedAt, p.ID, p.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	_, err := db.Exec("UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL", p.ID, p.UserID)
	if err != nil {
		return err
	}
	return nil
}

// Проверка, является ли пользователь владельцем поста
func (p *Post) IsOwner(userID int) bool {
	return p.UserID == userID
}

// Получить посты конкретного пользователя
func (p *Post) GetPostsByUser(db *sql.DB, userID int) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, user_id, title, content, created_at, updated_at, deleted_at FROM posts WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

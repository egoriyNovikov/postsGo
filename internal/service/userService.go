package service

import (
	"database/sql"

	"github.com/egoriyNovikov/internal/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	return user.CreateUser(s.db)
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	user := &models.User{ID: id}
	row := s.db.QueryRow("SELECT id, username, email, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL", id)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAllUsers() (*sql.Rows, error) {
	return s.db.Query("SELECT id, username, email, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL")
}

func (s *UserService) UpdateUser(user *models.User) error {
	_, err := s.db.Exec("UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4",
		user.Username, user.Email, user.UpdatedAt, user.ID)
	return err
}

func (s *UserService) DeleteUser(id int) error {
	_, err := s.db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

package service

import (
	"database/sql"

	"github.com/egoriyNovikov/internal/models"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		db: db,
	}
}

func (s *PostService) CreatePost(post *models.Post) error {
	return post.CreatePost(s.db)
}

func (s *PostService) GetPost(id int) (*models.Post, error) {
	post := &models.Post{ID: id}
	err := post.GetPost(s.db)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetAllPosts() (*sql.Rows, error) {
	post := &models.Post{}
	return post.GetAllPosts(s.db)
}

func (s *PostService) GetPostsByUser(userID int) (*sql.Rows, error) {
	post := &models.Post{}
	return post.GetPostsByUser(s.db, userID)
}

func (s *PostService) UpdatePost(id int, userID int, post *models.Post) error {
	post.ID = id
	post.UserID = userID
	return post.UpdatePost(s.db)
}

func (s *PostService) DeletePost(id int, userID int) error {
	post := &models.Post{ID: id, UserID: userID}
	return post.DeletePost(s.db)
}

func (s *PostService) IsPostOwner(postID int, userID int) (bool, error) {
	post, err := s.GetPost(postID)
	if err != nil {
		return false, err
	}
	return post.IsOwner(userID), nil
}

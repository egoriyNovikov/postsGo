package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/egoriyNovikov/internal/models"
	"github.com/egoriyNovikov/internal/service"
)

type PostController struct {
	postService *service.PostService
}

func NewPostController(postService *service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Получаем user_id из заголовка (установленного middleware)
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	post.UserID = userID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := c.postService.CreatePost(&post); err != nil {
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := c.postService.GetPost(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (c *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := c.postService.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var deletedAt sql.NullTime
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &deletedAt)
		if err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		if deletedAt.Valid {
			post.DeletedAt = &deletedAt.Time
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (c *PostController) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из заголовка (установленного middleware)
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	rows, err := c.postService.GetPostsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to get user posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var deletedAt sql.NullTime
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &deletedAt)
		if err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		if deletedAt.Valid {
			post.DeletedAt = &deletedAt.Time
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Получаем user_id из заголовка (установленного middleware)
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Проверяем, является ли пользователь владельцем поста
	isOwner, err := c.postService.IsPostOwner(id, userID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if !isOwner {
		http.Error(w, "Access denied: You can only edit your own posts", http.StatusForbidden)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	post.UpdatedAt = time.Now()

	if err := c.postService.UpdatePost(id, userID, &post); err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// Получаем обновленный пост
	updatedPost, err := c.postService.GetPost(id)
	if err != nil {
		http.Error(w, "Failed to get updated post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPost)
}

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Получаем user_id из заголовка (установленного middleware)
	userIDStr := r.Header.Get("User-ID")
	if userIDStr == "" {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Проверяем, является ли пользователь владельцем поста
	isOwner, err := c.postService.IsPostOwner(id, userID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if !isOwner {
		http.Error(w, "Access denied: You can only delete your own posts", http.StatusForbidden)
		return
	}

	if err := c.postService.DeletePost(id, userID); err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}

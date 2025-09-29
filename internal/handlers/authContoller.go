package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/egoriyNovikov/internal/models"
	"github.com/egoriyNovikov/pkg"
)

type AuthController struct {
	db *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		db: db,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var user models.User
	row := c.db.QueryRow("SELECT id, username, email, password_hash, created_at, updated_at, deleted_at FROM users WHERE username = $1 AND deleted_at IS NULL", loginReq.Username)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !pkg.CheckPasswordHash(loginReq.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := pkg.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken := pkg.GenerateRefreshToken()

	auth := models.Auth{
		User_id:    user.ID,
		Token:      refreshToken,
		Expires_at: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := auth.CreateAuth(c.db); err != nil {
		log.Printf("Error saving refresh token: %v", err)
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	response := LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		http.Error(w, "No refresh token provided", http.StatusBadRequest)
		return
	}

	_, err := c.db.Exec("DELETE FROM auth WHERE token = $1", refreshToken)
	if err != nil {
		log.Printf("Error deleting refresh token: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}

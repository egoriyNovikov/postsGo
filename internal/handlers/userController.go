package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/egoriyNovikov/internal/models"
	"github.com/egoriyNovikov/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := c.userService.CreateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := c.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

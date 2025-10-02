package router

import (
	"database/sql"
	"net/http"

	"github.com/egoriyNovikov/internal/handlers"
	"github.com/egoriyNovikov/internal/middleware"
	"github.com/egoriyNovikov/internal/service"
)

func Router(db *sql.DB, port string) {
	userService := service.NewUserService(db)
	postService := service.NewPostService(db)

	userController := handlers.NewUserController(userService)
	postController := handlers.NewPostController(postService)
	authController := handlers.NewAuthController(db)

	http.HandleFunc("/api/auth/login", authController.Login)
	http.HandleFunc("/api/auth/logout", middleware.RefreshTokenMiddleware(authController.Logout))

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userController.CreateUser(w, r)
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				middleware.AuthMiddleware(userController.GetUser)(w, r)
			} else {
				middleware.AuthMiddleware(userController.GetAllUsers)(w, r)
			}
		case http.MethodPut:
			middleware.AuthMiddleware(userController.UpdateUser)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(userController.DeleteUser)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(postController.CreatePost)(w, r)
		case http.MethodGet:
			// Проверяем, есть ли параметр id в запросе
			id := r.URL.Query().Get("id")
			if id != "" {
				// Если есть id, получаем конкретный пост
				middleware.AuthMiddleware(postController.GetPost)(w, r)
			} else {
				// Если нет id, получаем все посты
				middleware.AuthMiddleware(postController.GetAllPosts)(w, r)
			}
		case http.MethodPut:
			middleware.AuthMiddleware(postController.UpdatePost)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(postController.DeletePost)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/posts/my", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(postController.GetUserPosts)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+port, nil)
}

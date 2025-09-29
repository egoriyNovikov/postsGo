package router

import (
	"database/sql"
	"net/http"

	"github.com/egoriyNovikov/internal/handlers"
	"github.com/egoriyNovikov/internal/middleware"
	"github.com/egoriyNovikov/internal/service"
)

func Router(db *sql.DB, port string) {
	// Создаем сервисы
	userService := service.NewUserService(db)

	// Создаем контроллеры
	userController := handlers.NewUserController(userService)
	authController := handlers.NewAuthController(db)

	// Регистрируем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Маршруты аутентификации
	http.HandleFunc("/api/auth/login", authController.Login)
	http.HandleFunc("/api/auth/logout", middleware.RefreshTokenMiddleware(authController.Logout))

	// API маршруты для пользователей
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userController.CreateUser(w, r)
		case http.MethodGet:
			middleware.AuthMiddleware(userController.GetAllUsers)(w, r)
		case http.MethodPut:
			middleware.AuthMiddleware(userController.UpdateUser)(w, r)
		case http.MethodDelete:
			middleware.AuthMiddleware(userController.DeleteUser)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(userController.GetUser)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+port, nil)
}

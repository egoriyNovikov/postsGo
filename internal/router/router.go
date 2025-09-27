package router

import (
	"database/sql"
	"net/http"

	"github.com/egoriyNovikov/internal/handlers"
	"github.com/egoriyNovikov/internal/service"
)

func Router(db *sql.DB, port string) {
	// Создаем сервисы
	userService := service.NewUserService(db)

	// Создаем контроллеры
	userController := handlers.NewUserController(userService)

	// Регистрируем маршруты
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API маршруты для пользователей
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userController.CreateUser(w, r)
		case http.MethodGet:
			userController.GetAllUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userController.GetUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+port, nil)
}

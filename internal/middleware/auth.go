package middleware

import (
	"net/http"
	"strings"

	"github.com/egoriyNovikov/pkg"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := pkg.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("User-ID", string(rune(claims.UserID)))
		r.Header.Set("Username", claims.Username)

		next(w, r)
	}
}


func RefreshTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken := r.Header.Get("Authorization")
		if refreshToken == "" {
			http.Error(w, "Refresh token required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

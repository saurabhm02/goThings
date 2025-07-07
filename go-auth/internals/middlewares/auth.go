package middlewares

import (
	"context"
	"go-auth/internals/models"
	"go-auth/internals/utils"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserRoleKey contextKey = "userRole"

func RoleAuthorizationMiddleware(allowedRole models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			role, err := utils.GetRoleFromToken(token)
			if err != nil {
				log.Println("Invalid token:", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if role != allowedRole {
				http.Error(w, "Forbidden: You don't have access to this resource", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

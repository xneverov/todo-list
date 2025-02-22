package auth

import (
	"net/http"

	"github.com/xneverov/todo-list/internal/config"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		password := config.Get("TODO_PASSWORD")
		if password != "" {
			var token string

			cookie, err := req.Cookie("token")
			if err == nil {
				token = cookie.Value
			}

			valid, err := validateToken(token)
			if err != nil {
				valid = false
			}

			if !valid {
				http.Error(res, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(res, req)
	}
}

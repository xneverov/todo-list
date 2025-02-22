package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/xneverov/todo-list/internal/config"
)

type authResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandleAuth(res http.ResponseWriter, req *http.Request) {
	var authRequest struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(req.Body).Decode(&authRequest); err != nil {
		_ = json.NewEncoder(res).Encode(authResponse{Error: "Invalid JSON"})
		return
	}

	password := config.Get("TODO_PASSWORD")
	passwordEntered := authRequest.Password

	if passwordEntered != password {
		_ = json.NewEncoder(res).Encode(authResponse{Error: "Неверный пароль"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordEntered), bcrypt.DefaultCost)
	if err != nil {
		_ = json.NewEncoder(res).Encode(authResponse{Error: err.Error()})
		return
	}

	token, err := generateToken(string(passwordHash))
	if err != nil {
		_ = json.NewEncoder(res).Encode(authResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(res).Encode(authResponse{Token: token})
}

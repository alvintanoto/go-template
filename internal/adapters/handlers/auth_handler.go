package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"alvintanoto.id/go-template/internal/application"
)

type SignUpPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" || payload.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	err := h.authService.SignUp(r.Context(), email, payload.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "invalid email") {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}

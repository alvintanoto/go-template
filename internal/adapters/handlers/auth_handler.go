package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"alvintanoto.id/go-template/internal/application"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log         *zap.Logger
	authService *application.AuthService
}

func NewAuthHandler(log *zap.Logger, authService *application.AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}

type SignUpPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Error("decoding error", zap.Error(err))
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
		h.log.Error("sign up error", zap.Error(err))

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
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Error("decoding error", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" || payload.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.authService.ExecLogin(r.Context(), email, payload.Password)
	if err != nil {
		h.log.Error("login error", zap.Error(err))
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "login success",
		"token":   token,
	})
}

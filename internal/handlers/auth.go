package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/zostay/bobsbinder/internal/config"
)

type AuthHandler struct {
	DB     *sql.DB
	Config *config.Config
	Logger *zap.Logger
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, `{"error":"email, password, and name are required"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("failed to hash password", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?)",
		req.Email, string(hash), req.Name,
	)
	if err != nil {
		h.Logger.Error("failed to create user", zap.Error(err))
		http.Error(w, `{"error":"email already exists"}`, http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()

	// Create "self" party for the user
	_, err = h.DB.Exec(
		"INSERT INTO parties (user_id, name, relationship) VALUES (?, ?, 'self')",
		userID, req.Name,
	)
	if err != nil {
		h.Logger.Error("failed to create self party", zap.Error(err))
	}

	token, err := h.generateToken(userID, req.Email)
	if err != nil {
		h.Logger.Error("failed to generate token", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := authResponse{}
	resp.Token = token
	resp.User.ID = userID
	resp.User.Email = req.Email
	resp.User.Name = req.Name

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	var userID int64
	var email, name, passwordHash string
	err := h.DB.QueryRow(
		"SELECT id, email, name, password_hash FROM users WHERE email = ?",
		req.Email,
	).Scan(&userID, &email, &name, &passwordHash)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := h.generateToken(userID, email)
	if err != nil {
		h.Logger.Error("failed to generate token", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := authResponse{}
	resp.Token = token
	resp.User.ID = userID
	resp.User.Email = email
	resp.User.Name = name

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var email, name string
	err := h.DB.QueryRow("SELECT email, name FROM users WHERE id = ?", userID).Scan(&email, &name)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	token, err := h.generateToken(userID, email)
	if err != nil {
		h.Logger.Error("failed to generate token", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := authResponse{}
	resp.Token = token
	resp.User.ID = userID
	resp.User.Email = email
	resp.User.Name = name

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) generateToken(userID int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.Config.JWTSecret))
}

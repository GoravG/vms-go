package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"vms_go/internal/security"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB         *sql.DB
	HMACSecret []byte
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var userID int
	var passwordHash string
	err := h.DB.QueryRow("SELECT id, password_hash FROM users WHERE email=?", email).
		Scan(&userID, &passwordHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	data := []byte(email)

	tokenStr := security.GenerateHMAC(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Auth-Token", tokenStr)

	// Write JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"token":   tokenStr,
		"user_id": userID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		http.Error(w, "missing token", http.StatusUnauthorized)
// 		return
// 	}
// 	tokenStr := authHeader[len("Bearer "):]

// 	_, err := h.DB.Exec("DELETE FROM active_tokens WHERE token=$1", tokenStr)
// 	if err != nil {
// 		http.Error(w, "could not logout", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte("logged out"))
// }

// func (h *AuthHandler) Protected(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("You have accessed a protected route."))
// }

package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"vms_go/internal/token"
	"vms_go/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if password == "" || email == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(email) {
		http.Error(w, "Email is not valid", http.StatusBadRequest)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "could not hash password", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(`INSERT INTO users (email, password_hash) VALUES (?, ?)`, email, string(hash))
	if err != nil {
		log.Printf("DB insert error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
}

func (h *UserHandler) Checkin(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	paramToken := params.Get("token")
	utils.LogInfo(paramToken)
	if paramToken == "" {
		http.Error(w, "missing required parameters: token", http.StatusBadRequest)
		return
	}
	redisToken := token.GetToken()
	if paramToken == redisToken {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

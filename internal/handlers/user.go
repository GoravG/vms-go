package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"vms_go/internal/security"
	"vms_go/internal/token"
	"vms_go/internal/utils"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := utils.CreateUser(h.DB, email, password)
	if err != nil {
		log.Printf("Create user error: %v", err)
		statusCode := http.StatusInternalServerError

		// Set appropriate status code based on error type
		switch err.Error() {
		case "email and password are required", "email is not valid":
			statusCode = http.StatusBadRequest
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created"))
}

func (h *UserHandler) Checkin(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	paramToken := params.Get("token")
	if paramToken == "" {
		http.Error(w, "missing required parameter: token", http.StatusBadRequest)
		return
	}

	var body struct {
		UserToken string `json:"user_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if body.UserToken == "" {
		http.Error(w, "missing required field: user_token", http.StatusBadRequest)
		return
	}

	claims, err := security.VerifyAndExtractClaims(body.UserToken)
	if err != nil {
		http.Error(w, "unable to verify token", http.StatusUnauthorized)
	}

	// Compare token with Redis value (as before)
	redisToken := token.GetToken()
	if paramToken == redisToken {
		w.WriteHeader(http.StatusOK)
		alreadyCheckedIn, err := utils.VisitorAlreadyCheckedIn(h.DB, claims.Email)
		if err != nil {
			log.Fatal(err)
		}
		if alreadyCheckedIn {
			err := utils.UpdateVisitLog(h.DB, claims.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error in updating exising visit log"))
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Updated visit log successfully"))
		} else {
			err := utils.InsertVisitLog(h.DB, claims.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error in inserting visit log"))
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Checked in successfully"))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

package utils

import (
	"database/sql"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return matched
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, email, password string) error {
	if password == "" || email == "" {
		return errors.New("email and password are required")
	}

	if !IsValidEmail(email) {
		return errors.New("email is not valid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}

	const query = `
		INSERT INTO users (email, password_hash)
		VALUES (?, ?)
	`

	_, err = db.Exec(query, email, string(hash))
	if err != nil {
		return err
	}

	return nil
}

func VisitorAlreadyCheckedIn(db *sql.DB, email string) (bool, error) {
	const query = `
	SELECT EXISTS(SELECT 1 FROM visit_log WHERE
	visitor = ? AND visit_date = CURDATE() AND check_in_time IS NOT NULL
	)`
	var alreadyCheckedIn bool
	err := db.QueryRow(query, email).Scan(&alreadyCheckedIn)
	return alreadyCheckedIn, err
}

func UpdateVisitLog(db *sql.DB, email string) error {
	const query = `
	UPDATE visit_log SET check_out_time = NOW()
	WHERE visitor = ? AND visit_date = CURDATE()
	`
	_, err := db.Exec(query, email)
	return err
}

func InsertVisitLog(db *sql.DB, email string) error {
	const query = `
	INSERT INTO visit_log (visitor, visit_date, check_in_time, check_out_time)
	VALUES (?, CURDATE(), NOW(), NULL)
	`
	_, err := db.Exec(query, email)
	return err
}

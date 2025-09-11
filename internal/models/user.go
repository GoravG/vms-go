package models

import "time"

type User struct {
	ID           int64
	PasswordHash string
	Email        string
	FirstName    string
	LastName     string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	ModifiedAt   time.Time
}

func (User) CreateTableSQL() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		password_hash VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		role ENUM('admin', 'user') DEFAULT 'user',
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
}

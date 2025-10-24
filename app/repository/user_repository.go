package repository

import (
	"database/sql"
	"prak4/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	err := r.DB.QueryRow(
		"SELECT id, username, email, password_hash, role, created_at FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(username, email, passwordHash, role string) (*model.User, error) {
	var user model.User

	// SQL query to insert new user
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING id, username, email, role, created_at`

	err := r.DB.QueryRow(
		query,
		username,
		email,
		passwordHash,
		role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
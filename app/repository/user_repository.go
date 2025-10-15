package repository

import (
	"database/sql"
	"prak4/app/model"
)

type UserRepository struct {
	DB *sql.DB
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